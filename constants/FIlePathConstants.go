package constants

import (
	"path/filepath"
)

const AppFolder string = "Clave"

var AppFolderPATH string = filepath.Join(ProgramData, AppFolder)

var DatabaseLocation string = filepath.Join(AppFolderPATH, "Databases")

var SecureVaultDB string = filepath.Join(DatabaseLocation, "secure_vault.db")
