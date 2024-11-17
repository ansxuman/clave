package backend

import (
	"clave/constants"
	"clave/localstorage"
	"clave/objects"
	"clave/services/auth"
	"clave/services/window"
	"runtime"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	cancel      cmap.ConcurrentMap[string, func()]
	authService *auth.Service
	winManager  *window.Manager
	isVerified  bool
	isMacOS     bool
}

func NewApp() *App {
	storage := localstorage.GetPersistentStorage()
	authService := auth.NewService(storage)

	app := &App{
		cancel:      cmap.New[func()](),
		authService: authService,
		isVerified:  false,
		isMacOS:     runtime.GOOS == "darwin",
	}

	app.winManager = window.NewManager(app)
	return app
}

func (a *App) HasPin() bool {
	return a.authService.HasPin()
}

func (a *App) IsVerified() bool {
	return a.isVerified
}

func (a *App) SetVerified(state bool) {
	a.isVerified = state
}

func (a *App) Initialize() objects.InitResult {
	if !a.HasPin() {
		return objects.InitResult{NeedsOnboarding: true}
	}
	return objects.InitResult{NeedsVerification: false}
}

func (a *App) SetWindow(window *application.WebviewWindow) {
	a.winManager.SetWindow(window)
}

func (a *App) SetupPin(pin string) error {
	return a.authService.SetupPin(pin)
}

func (a *App) VerifyPin(pin string) (bool, error) {
	isValid, err := a.authService.VerifyPin(pin)
	if isValid {
		a.SetVerified(true)
	}
	return isValid, err
}

func (a *App) IsMacOS() bool {
	return a.isMacOS
}

func (a *App) VerifyTouchID() bool {
	return a.winManager.HandleTouchID()
}

func (a *App) GetAppVersion() string {
	return constants.AppVersion
}
