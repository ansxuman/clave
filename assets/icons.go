package constants

import (
	_ "embed"
)

//go:embed password.png
var iconData []byte

//go:embed passworddark.png
var darkIconData []byte

// Icon for the application
var Icon = iconData
var DarkIcon = darkIconData
