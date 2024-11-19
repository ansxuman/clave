package totp

import (
	"bytes"
	"clave/objects"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/liyue201/goqr"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	storage Storage
	window  *application.WebviewWindow
}

type Storage interface {
	CheckIfIssuerOrSecretExists(issuer, secret string) (bool, objects.TotpSecretObject)
	AddTotpSecretObject(issuer, secret string) error
	GetListOfTotpSecretObjects() objects.TotpSecretObjectList
	DeleteTotpSecretObject(profileId string) error
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) SetWindow(window *application.WebviewWindow) {
	if window != nil {
		s.window = window
	}
}

func (s *Service) OpenQR() error {
	dialog := application.OpenFileDialog().
		SetTitle("Select QR Code Image").
		AddFilter("Image Files", "*.png;*.jpg;*.jpeg")

	result, err := dialog.PromptForSingleSelection()
	if err != nil {
		log.Printf("[TOTP] Failed to open QR image: %v", err)
		return err
	}

	if result != "" {
		return s.scanQRCode(result)
	}
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

func (s *Service) configureQRProfile(qrData []uint8) error {
	parsedURL, err := url.Parse(string(qrData))
	if err != nil {
		log.Printf("[TOTP] Failed to parse QR data: %v", err)
		s.window.EmitEvent("failedToScanQR", "Invalid QR code format")
		return err
	}

	if parsedURL.Scheme != "otpauth" || parsedURL.Host != "totp" {
		s.window.EmitEvent("failedToScanQR", "Invalid TOTP URI format")
		return fmt.Errorf("invalid URI format")
	}

	secret := parsedURL.Query().Get("secret")
	if secret == "" {
		s.window.EmitEvent("failedToScanQR", "Missing secret parameter")
		return fmt.Errorf("missing secret")
	}
	base32Regex := regexp.MustCompile(`^[A-Z2-7]+=*$`)
	if !base32Regex.MatchString(secret) {
		s.window.EmitEvent("failedToScanQR", "Invalid secret format")
		return fmt.Errorf("invalid secret format")
	}

	path := strings.TrimPrefix(parsedURL.Path, "/")
	issuer := strings.TrimSuffix(path, "/")
	if issuer == "" {
		s.window.EmitEvent("failedToScanQR", "Missing issuer")
		return fmt.Errorf("missing issuer")
	}

	exists, profile := s.storage.CheckIfIssuerOrSecretExists(issuer, secret)
	if exists {
		log.Printf("[TOTP] Profile already exists for issuer: %s", issuer)
		s.window.EmitEvent("duplicateScanQR", profile)
		return nil
	}

	if err := s.storage.AddTotpSecretObject(issuer, secret); err != nil {
		log.Printf("[TOTP] Failed to save profile: %v", err)
		return err
	}

	s.window.EmitEvent("refreshTOTPProfiles", nil)
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
		s.window.EmitEvent("totpData", []objects.TotpSecretObject{})
		return
	}

	log.Printf("[TOTP] Sending %d profiles", len(listOfSecrets))
	s.window.EmitEvent("totpData", listOfSecrets)
}

func (s *Service) RemoveTotpProfile(profileId string) error {
	return s.storage.DeleteTotpSecretObject(profileId)
}
