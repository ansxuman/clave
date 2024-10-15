package backend

var (
	BuildVersion = "0.0.1"
	Version      = BuildVersion
)

func (a *App) GetAppVersion() string {
	return Version
}
