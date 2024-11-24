package window

import (
	"github.com/ansxuman/go-touchid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Manager struct {
	window          *application.WebviewWindow
	authService     AuthService
	isAddingProfile bool
}

type AuthService interface {
	HasPin() bool
	IsVerified() bool
	SetVerified(bool)
	VerifyTouchID() bool
}

func NewManager(authService AuthService) *Manager {
	return &Manager{
		authService:     authService,
		isAddingProfile: false,
	}
}

func (m *Manager) SetWindow(window *application.WebviewWindow) {
	m.window = window
	m.setupWindowEvents()
}

func (m *Manager) setupWindowEvents() {
	focusHandler := func(e *application.WindowEvent) {
		if m.isAddingProfile {
			return
		}

		if !m.authService.HasPin() || m.authService.IsVerified() {
			return
		}

		m.window.EmitEvent("requirePinVerification", nil)
	}

	m.window.OnWindowEvent(events.Common.WindowFocus, focusHandler)
	m.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
		if !m.isAddingProfile {
			m.authService.SetVerified(false)
		}
	})
}

func (m *Manager) HandleTouchID() bool {
	success, err := touchid.Auth(touchid.DeviceTypeBiometrics, "Verify Identity")
	if err == nil && success {
		m.authService.SetVerified(true)
		m.window.EmitEvent("verificationComplete", nil)
		return true
	}
	return false
}

func (m *Manager) StartProfileAddition() {
	m.isAddingProfile = true
}

func (m *Manager) EndProfileAddition() {
	m.isAddingProfile = false
}
