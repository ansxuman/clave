package main

import (
	a "clave/assets"
	"clave/backend"
	"clave/constants"
	"embed"
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	backendApp := backend.NewApp()
	app := application.New(application.Options{
		Name:        constants.ApplicationName,
		Description: constants.Description,
		Services: []application.Service{
			application.NewService(backendApp, application.ServiceOptions{}),
		},
		OnShutdown: func() {
			log.Println("Shutting down...")
		},
		ShouldQuit: func() bool {
			return true
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systemTray := app.NewSystemTray()

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:     constants.MaxWidth,
		Height:    constants.MaxHeight,
		Title:     constants.ApplicationName,
		Frameless: true,
		// Always on top is disabled
		// because its crashing the
		// application
		// AlwaysOnTop:       true,
		Hidden:          true,
		DisableResize:   true,
		DevToolsEnabled: false,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			Appearance:              application.NSAppearanceNameAccessibilityHighContrastVibrantLight,
		},
	})

	app.ShowAboutDialog()

	backendApp.SetWindow(window)

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		log.Println("Application started, initializing security...")
	})

	app.Hide()

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(a.Icon)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Clave").SetEnabled(false)

	myMenu.Add("Go Back to App").OnClick(func(ctx *application.Context) {
		window.Show()
	})

	myMenu.AddSeparator()
	myMenu.Add("Quit Clave").OnClick(func(ctx *application.Context) {
		q := application.QuestionDialog().
			SetTitle("Quit Clave").
			SetMessage("Are you sure you want to quit?")

		q.AddButton("Yes").OnClick(func() {
			app.Quit()
		})

		q.AddButton("No").SetAsDefault().OnClick(func() {

		})

		q.Show()
	})

	systemTray.SetMenu(myMenu)

	systemTray.AttachWindow(window).WindowOffset(5)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
