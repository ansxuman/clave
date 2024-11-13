package constants

import (
	"os"
)

var homeDir, _ = os.UserHomeDir()
var ProgramData = homeDir
