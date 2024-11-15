package backend

import (
	"clave/constants"
	"clave/localstorage"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"runtime"

	"github.com/ansxuman/go-touchid"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"golang.org/x/crypto/argon2"
)

const (
	PIN_KEY  = "user_pin"
	SALT_KEY = "pin_salt"
)

type App struct {
	ctx        context.Context
	cancel     cmap.ConcurrentMap[string, func()]
	storage    *localstorage.PersistentStore
	isVerified bool
	window     *application.WebviewWindow
	isMacOS    bool
}

func NewApp() *App {
	return &App{
		cancel:     cmap.New[func()](),
		storage:    localstorage.GetPersistentStorage(),
		isVerified: false,
		isMacOS:    runtime.GOOS == "darwin",
	}
}

type InitResult struct {
	NeedsOnboarding   bool `json:"needsOnboarding"`
	NeedsVerification bool `json:"needsVerification"`
}

func (a *App) Initialize() InitResult {
	a.window.EmitEvent("appVersion", constants.AppVersion)
	if !a.storage.IsHealthy() {
		panic("Storage is not healthy")
	}

	if !a.HasPin() {
		return InitResult{NeedsOnboarding: true}
	}

	return InitResult{NeedsVerification: false}
}

func (a *App) HasPin() bool {
	return a.storage.HasKey(PIN_KEY) && a.storage.HasKey(SALT_KEY)
}

func (a *App) SetupPin(pin string) error {
	if len(pin) < 4 {
		return errors.New("PIN must be at least 4 characters")
	}
	if a.HasPin() {
		return errors.New("PIN is already set")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(pin), salt, 1, 64*1024, 4, 32)

	if err := a.storage.SetValue(SALT_KEY, base64.StdEncoding.EncodeToString(salt)); err != nil {
		return err
	}

	return a.storage.SetValue(PIN_KEY, base64.StdEncoding.EncodeToString(hash))
}

func (a *App) SetVerified(state bool) {
	a.isVerified = state
}

func (a *App) IsVerified() bool {
	return a.isVerified
}

func (a *App) VerifyPin(pin string) (bool, error) {
	if !a.storage.IsHealthy() {
		return false, errors.New("storage is not healthy")
	}
	var storedHashStr string
	var saltStr string

	if err := a.storage.Get(PIN_KEY, &storedHashStr); err != nil {
		return false, err
	}

	if err := a.storage.Get(SALT_KEY, &saltStr); err != nil {
		return false, err
	}

	storedHash, err := base64.StdEncoding.DecodeString(storedHashStr)
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(pin), salt, 1, 64*1024, 4, 32)

	isValid := subtle.ConstantTimeCompare(hash, storedHash) == 1
	if isValid {
		a.SetVerified(true)
	}
	return isValid, nil
}

func (a *App) SetWindow(window *application.WebviewWindow) {
	a.window = window

	focusHandler := func(e *application.WindowEvent) {
		if !a.HasPin() || a.isVerified {
			return
		}

		if runtime.GOOS == "darwin" {
			success, err := touchid.Auth(touchid.DeviceTypeBiometrics, "Verify Identity")
			if err == nil && success {
				a.isVerified = true
				window.EmitEvent("verificationComplete", nil)
				return
			}
		}

		window.EmitEvent("requirePinVerification", nil)
	}

	window.OnWindowEvent(events.Common.WindowFocus, focusHandler)
	window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
		a.isVerified = false
	})
}
