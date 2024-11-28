package constants

import (
	"os"
	"path/filepath"
)

var homeDir, _ = os.UserHomeDir()
var ProgramData = filepath.Join(homeDir, "Library", "Application Support")
