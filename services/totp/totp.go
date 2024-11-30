package totp

import (
	"bytes"
	"clave/localstorage"
	"clave/objects"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/liyue201/goqr"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	storage    Storage
	window     *application.WebviewWindow
	winManager WindowManager
}

type Storage interface {
	CheckIfIssuerOrSecretExists(issuer, secret string) (bool, objects.TotpSecretObject)
	AddTotpSecretObject(issuer, secret string) error
	GetListOfTotpSecretObjects() objects.TotpSecretObjectList
	DeleteTotpSecretObject(profileId string) error
	SetValue(key string, value interface{}) error
	Get(key string, value interface{}) error
	DeleteKey(key string) error
}

type WindowManager interface {
	StartProfileAddition()
	EndProfileAddition()
}

func NewService(storage Storage, winManager WindowManager) *Service {
	return &Service{
		storage:    storage,
		winManager: winManager,
	}
}

func (s *Service) SetWindow(window *application.WebviewWindow) {
	if window != nil {
		s.window = window
	}
}

func (s *Service) OpenQR() error {
	s.winManager.StartProfileAddition()

	dialog := application.OpenFileDialog().
		SetTitle("Select QR Code Image").
		AddFilter("Image Files", "*.png;*.jpg;*.jpeg")

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		s.winManager.EndProfileAddition()
		log.Printf("[TOTP] Failed to open QR image: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventQRScanError, []string{"Unable to open image selector"})
		return err
	}

	if result != "" {
		err = s.scanQRCode(result)
		s.winManager.EndProfileAddition()
		return err
	}

	s.winManager.EndProfileAddition()
	return nil
}

func (s *Service) scanQRCode(path string) error {
	imgdata, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[TOTP] Failed to read image: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventQRScanError, []string{"Unable to read the selected image"})
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Printf("[TOTP] Failed to decode image: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventQRScanError, []string{"The selected file is not a valid image"})
		return err
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Printf("[TOTP] Failed to recognize QR: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventQRScanError, []string{"No QR code found in the image"})
		return err
	}

	if len(qrCodes) == 0 {
		s.window.Show()
		s.window.EmitEvent(EventQRScanError, []string{"No valid QR code found in the image"})
		return fmt.Errorf("no QR codes found")
	}

	for _, qrCode := range qrCodes {
		if err := s.configureQRProfile(qrCode.Payload); err != nil {
			return err
		}
	}
	return nil
}

const (
	EventTOTPData           = "totpData"
	EventFailedToAddProfile = "failedToAddProfile"
	EventRefreshProfiles    = "refreshProfiles"
	EventDuplicateProfile   = "duplicateProfile"
	EventBackupError        = "backupError"
	EventBackupSuccess      = "backupSuccess"
	EventRestoreError       = "restoreError"
	EventRestoreSuccess     = "restoreSuccess"
	EventQRScanError        = "qrScanError"
)

func (s *Service) configureQRProfile(qrData []uint8) error {
	parsedURL, err := url.Parse(string(qrData))
	if err != nil {
		log.Printf("[TOTP] Failed to parse QR data: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Invalid QR code format")
		return err
	}

	if parsedURL.Scheme != "otpauth" || parsedURL.Host != "totp" {
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Invalid QR code")
		return fmt.Errorf("invalid URI format")
	}

	secret := parsedURL.Query().Get("secret")
	if secret == "" {
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Missing secret parameter")
		return fmt.Errorf("missing secret")
	}
	base32Regex := regexp.MustCompile(`^[A-Z2-7]+=*$`)
	if !base32Regex.MatchString(secret) {
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Invalid secret format")
		return fmt.Errorf("invalid secret format")
	}

	path := strings.TrimPrefix(parsedURL.Path, "/")
	issuer := strings.TrimSuffix(path, "/")
	if issuer == "" {
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Missing issuer")
		return fmt.Errorf("missing issuer")
	}

	exists, profile := s.storage.CheckIfIssuerOrSecretExists(issuer, secret)
	if exists {
		s.window.Show()
		log.Printf("[TOTP] Profile already exists for issuer: %s", issuer)
		s.window.EmitEvent(EventDuplicateProfile, fmt.Sprintf("Profile '%s' already exists", profile.Issuer))
		return nil
	}

	if err := s.storage.AddTotpSecretObject(issuer, secret); err != nil {
		log.Printf("[TOTP] Failed to save profile: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventFailedToAddProfile, "Failed to save profile")
		s.winManager.EndProfileAddition()
		return err
	}
	s.window.Show()
	s.window.EmitEvent(EventRefreshProfiles, nil)
	s.winManager.EndProfileAddition()
	return nil
}

func (s *Service) SendTOTPData() {
	if s.window == nil {
		log.Printf("[TOTP] Window not initialized")
		return
	}

	listOfSecrets := s.storage.GetListOfTotpSecretObjects()
	if len(listOfSecrets) == 0 {
		log.Printf("[TOTP] No profiles found")
		s.window.EmitEvent(EventTOTPData, []objects.TotpSecretObject{})
		return
	}

	log.Printf("[TOTP] Sending %d profiles", len(listOfSecrets))
	s.window.EmitEvent(EventTOTPData, listOfSecrets)
}

func (s *Service) RemoveTotpProfile(profileId string) error {
	return s.storage.DeleteTotpSecretObject(profileId)
}

func (s *Service) AddManualProfile(issuer string, secret string) error {
	base32Regex := regexp.MustCompile(`^[A-Z2-7]+=*$`)
	if !base32Regex.MatchString(secret) {
		s.window.EmitEvent(EventFailedToAddProfile, "Invalid secret format")
		return fmt.Errorf("invalid secret format")
	}

	exists, profile := s.storage.CheckIfIssuerOrSecretExists(issuer, secret)
	if exists {
		log.Printf("[TOTP] Profile already exists for issuer: %s", issuer)
		s.window.EmitEvent(EventDuplicateProfile, fmt.Sprintf("Profile '%s' already exists", profile.Issuer))
		return nil
	}

	if err := s.storage.AddTotpSecretObject(issuer, secret); err != nil {
		log.Printf("[TOTP] Failed to save profile: %v", err)
		s.window.EmitEvent(EventFailedToAddProfile, "Failed to add profile")
		return err
	}

	s.window.EmitEvent("refreshTOTPProfiles", nil)
	return nil
}

func (s *Service) BackupProfiles() error {
	log.Printf("[TOTP] Starting backup process")
	s.winManager.StartProfileAddition()

	existingProfiles := s.storage.GetListOfTotpSecretObjects()
	if len(existingProfiles) == 0 {
		log.Printf("[TOTP] No profiles found to backup")
		s.window.Show()
		s.window.EmitEvent(EventBackupError, []string{"You don't have any profiles to backup yet"})
		s.winManager.EndProfileAddition()
		return fmt.Errorf("no profiles to backup")
	}

	dialog := application.SaveFileDialog().
		SetMessage("Save backup file").
		SetButtonText("Save Backup").
		AddFilter("Backup Files", "*.clave").
		SetFilename(fmt.Sprintf("clave-backup-%s.clave", time.Now().Format("2006-01-02-15-04-05")))

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		log.Printf("[TOTP] Backup dialog error: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventBackupError, []string{"Unable to open save dialog"})
		s.winManager.EndProfileAddition()
		return err
	}

	if result == "" {
		log.Printf("[TOTP] User cancelled backup")
		s.winManager.EndProfileAddition()
		return nil
	}

	backupData := struct {
		Profiles []objects.TotpSecretObject `json:"profiles"`
		Version  string                     `json:"version"`
		Date     string                     `json:"date"`
	}{
		Profiles: existingProfiles,
		Version:  "1.0",
		Date:     time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		log.Printf("[TOTP] Failed to prepare backup data: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventBackupError, []string{"Unable to prepare backup data"})
		s.winManager.EndProfileAddition()
		return err
	}

	encrypted, err := s.storage.(*localstorage.PersistentStore).Encrypt(jsonData)
	if err != nil {
		log.Printf("[TOTP] Failed to encrypt backup: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventBackupError, []string{"Unable to secure backup data"})
		s.winManager.EndProfileAddition()
		return err
	}

	err = os.WriteFile(result, encrypted, 0644)
	if err != nil {
		log.Printf("[TOTP] Failed to write backup file: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventBackupError, []string{"Unable to save backup file"})
		s.winManager.EndProfileAddition()
		return err
	}

	s.window.Show()
	s.window.EmitEvent(EventBackupSuccess, []string{fmt.Sprintf("Successfully backed up %d profiles", len(existingProfiles))})
	s.winManager.EndProfileAddition()
	return nil
}

func (s *Service) RestoreProfiles() error {
	log.Printf("[TOTP] Starting restore process")
	s.winManager.StartProfileAddition()

	dialog := application.OpenFileDialog().
		SetTitle("Select Backup File").
		AddFilter("Backup Files", "*.clave")

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		log.Printf("[TOTP] Restore dialog error: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventRestoreError, []string{"Unable to open file selector"})
		s.winManager.EndProfileAddition()
		return err
	}

	if result == "" {
		log.Printf("[TOTP] User cancelled restore")
		s.winManager.EndProfileAddition()
		return nil
	}

	encryptedData, err := os.ReadFile(result)
	if err != nil {
		log.Printf("[TOTP] Failed to read backup file: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventRestoreError, []string{"Unable to read the backup file"})
		s.winManager.EndProfileAddition()
		return err
	}

	decrypted, err := s.storage.(*localstorage.PersistentStore).Decrypt(encryptedData)
	if err != nil {
		log.Printf("[TOTP] Failed to decrypt backup: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventRestoreError, []string{"Invalid or corrupted backup file"})
		s.winManager.EndProfileAddition()
		return err
	}

	var backupData struct {
		Profiles []objects.TotpSecretObject `json:"profiles"`
		Version  string                     `json:"version"`
		Date     string                     `json:"date"`
	}

	if err := json.Unmarshal(decrypted, &backupData); err != nil {
		log.Printf("[TOTP] Failed to parse backup data: %v", err)
		s.window.Show()
		s.window.EmitEvent(EventRestoreError, []string{"The backup file appears to be damaged"})
		s.winManager.EndProfileAddition()
		return err
	}

	if len(backupData.Profiles) == 0 {
		log.Printf("[TOTP] Empty backup file")
		s.window.Show()
		s.window.EmitEvent(EventRestoreError, []string{"No profiles found in the backup file"})
		s.winManager.EndProfileAddition()
		return fmt.Errorf("empty backup")
	}

	existingProfiles := s.storage.GetListOfTotpSecretObjects()
	existingMap := make(map[string]bool)
	for _, p := range existingProfiles {
		existingMap[p.Secret] = true
	}

	var stats struct {
		added     int
		duplicate int
		failed    int
	}

	for _, profile := range backupData.Profiles {
		if _, exists := existingMap[profile.Secret]; exists {
			log.Printf("[TOTP] Skipping duplicate profile: %s", profile.Issuer)
			stats.duplicate++
			continue
		}

		err := s.storage.AddTotpSecretObject(profile.Issuer, profile.Secret)
		if err != nil {
			log.Printf("[TOTP] Failed to restore profile %s: %v", profile.Issuer, err)
			stats.failed++
		} else {
			log.Printf("[TOTP] Added profile: %s", profile.Issuer)
			stats.added++
		}
	}

	var statusParts []string
	if stats.added > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d new profiles added", stats.added))
	}
	if stats.duplicate > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d profiles already existed", stats.duplicate))
	}
	if stats.failed > 0 {
		statusParts = append(statusParts, fmt.Sprintf("%d profiles couldn't be added", stats.failed))
	}

	statusMsg := strings.Join(statusParts, ", ")
	if stats.added == 0 && stats.duplicate > 0 {
		statusMsg = "All profiles from the backup already exist"
	}

	s.window.Show()
	s.window.EmitEvent(EventRestoreSuccess, []string{statusMsg})
	s.SendTOTPData()
	s.winManager.EndProfileAddition()
	return nil
}
