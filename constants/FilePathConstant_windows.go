package constants

import (
	"os"
)

var ProgramData string = os.Getenv("USERPROFILE")
