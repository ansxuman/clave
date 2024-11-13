package constants

import (
	_ "embed"
)

//go:embed password.png
var iconData []byte

// Icon for the application
var Icon = iconData
