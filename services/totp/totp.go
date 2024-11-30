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
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Printf("[TOTP] Failed to decode image: %v", err)
		s.window.EmitEvent("failedToScanQR", nil)
		return err
	}

	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Printf("[TOTP] Failed to recognize QR: %v", err)
		s.window.EmitEvent("failedToScanQR", nil)
		return err
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
		s.window.EmitEvent(EventFailedToAddProfile, "Invalid TOTP URI format")
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

	existingProfiles := s.storage.GetListOfTotpSecretObjects()
	if len(existingProfiles) == 0 {
		log.Printf("[TOTP] No profiles found to backup")
		s.window.EmitEvent("backupError", []string{"No profiles available to backup"})
		return fmt.Errorf("no profiles to backup")
	}

	dialog := application.SaveFileDialog().
		SetMessage("Save backup file").
		SetButtonText("Save Backup").
		AddFilter("Backup Files", "*.clave").
		SetFilename(fmt.Sprintf("clave-backup-%s.clave", time.Now().Format("2006-01-02-15-04-05")))

	log.Printf("[TOTP] Showing save dialog")
	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		log.Printf("[TOTP] Backup dialog error: %v", err)
		s.window.EmitEvent("backupError", []string{"Failed to open save dialog"})
		return err
	}

	if result == "" {
		log.Printf("[TOTP] User cancelled backup")
		return nil
	}
	log.Printf("[TOTP] Selected backup path: %s", result)

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
		log.Printf("[TOTP] Failed to marshal backup data: %v", err)
		s.window.EmitEvent("backupError", []string{"Failed to prepare backup data"})
		return err
	}

	encrypted, err := s.storage.(*localstorage.PersistentStore).Encrypt(jsonData)
	if err != nil {
		log.Printf("[TOTP] Failed to encrypt backup data: %v", err)
		s.window.EmitEvent("backupError", []string{"Failed to encrypt backup data"})
		return err
	}

	err = os.WriteFile(result, encrypted, 0644)
	if err != nil {
		log.Printf("[TOTP] Failed to write backup file: %v", err)
		s.window.EmitEvent("backupError", []string{"Failed to save backup file"})
		return err
	}

	log.Printf("[TOTP] Backup completed successfully")
	s.window.EmitEvent("backupSuccess", []string{fmt.Sprintf("Successfully backed up %d profiles", len(existingProfiles))})
	return nil
}

func (s *Service) RestoreProfiles() error {
	log.Printf("[TOTP] Starting restore process")

	dialog := application.OpenFileDialog().
		SetTitle("Select Backup File").
		AddFilter("Backup Files", "*.clave")

	log.Printf("[TOTP] Showing open dialog")
	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		log.Printf("[TOTP] Restore dialog error: %v", err)
		s.window.EmitEvent("restoreError", []string{"Failed to open file dialog"})
		return err
	}

	if result == "" {
		log.Printf("[TOTP] User cancelled restore")
		return nil
	}
	log.Printf("[TOTP] Selected restore file: %s", result)

	encryptedData, err := os.ReadFile(result)
	if err != nil {
		log.Printf("[TOTP] Failed to read backup file: %v", err)
		s.window.EmitEvent("restoreError", []string{"Failed to read backup file"})
		return err
	}

	decrypted, err := s.storage.(*localstorage.PersistentStore).Decrypt(encryptedData)
	if err != nil {
		log.Printf("[TOTP] Failed to decrypt backup file: %v", err)
		s.window.EmitEvent("restoreError", []string{"Failed to decrypt backup file"})
		return err
	}

	var backupData struct {
		Profiles []objects.TotpSecretObject `json:"profiles"`
		Version  string                     `json:"version"`
		Date     string                     `json:"date"`
	}

	if err := json.Unmarshal(decrypted, &backupData); err != nil {
		log.Printf("[TOTP] Failed to parse backup data: %v", err)
		s.window.EmitEvent("restoreError", []string{"Invalid backup file format"})
		return err
	}

	if len(backupData.Profiles) == 0 {
		log.Printf("[TOTP] No profiles found in backup file")
		s.window.EmitEvent("restoreError", []string{"No profiles found in backup"})
		return fmt.Errorf("no profiles in backup")
	}

	log.Printf("[TOTP] Found %d profiles in backup (Version: %s, Date: %s)",
		len(backupData.Profiles), backupData.Version, backupData.Date)

	existingProfiles := s.storage.GetListOfTotpSecretObjects()
	existingMap := make(map[string]bool)
	for _, p := range existingProfiles {
		existingMap[p.Secret] = true
	}

	var stats struct {
		added     int
		skipped   int
		duplicate int
	}

	for _, profile := range backupData.Profiles {
		if _, exists := existingMap[profile.Secret]; exists {
			log.Printf("[TOTP] Skipping duplicate profile: %s", profile.Issuer)
			stats.duplicate++
			continue
		}

		err := s.storage.AddTotpSecretObject(profile.Issuer, profile.Secret)
		if err != nil {
			log.Printf("[TOTP] Failed to add profile %s: %v", profile.Issuer, err)
			stats.skipped++
		} else {
			log.Printf("[TOTP] Added new profile: %s", profile.Issuer)
			stats.added++
		}
	}

	statusMsg := fmt.Sprintf("Restore completed: %d added", stats.added)
	if stats.duplicate > 0 {
		statusMsg += fmt.Sprintf(", %d already existed", stats.duplicate)
	}
	if stats.skipped > 0 {
		statusMsg += fmt.Sprintf(", %d failed", stats.skipped)
	}

	log.Printf("[TOTP] %s", statusMsg)
	s.window.EmitEvent("restoreSuccess", []string{statusMsg})
	s.SendTOTPData()
	return nil
}
