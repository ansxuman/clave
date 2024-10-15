package main

import (
	"clave/backend"
	"clave/constants"
	"embed"
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	backendApp := backend.NewApp()
	app := application.New(application.Options{
		Name:        "Clave",
		Description: "Secure Authentication at Your Fingertips",
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
	})

	systemTray := app.NewSystemTray()

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         350,
		Height:        500,
		Title:         "Clave",
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
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

	app.Hide()

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(constants.Icon)
	}

	myMenu := app.NewMenu()
	myMenu.Add("Clave").SetEnabled(false)

	myMenu.Add("Go Back to App").OnClick(func(ctx *application.Context) {
		window.Show()
	})

	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
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
