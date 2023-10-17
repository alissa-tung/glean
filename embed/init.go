package embed

import (
	_ "embed"
	"runtime"
)

//go:embed elan-init.ps1
var InitWindows []byte

//go:embed elan-init.sh
var InitUnix []byte

func InitScriptName() string {
	switch runtime.GOOS {
	case "windows":
		return "elan-init.ps1"
	default:
		return "elan-init.sh"
	}
}

func InitScriptBytes() []byte {
	switch runtime.GOOS {
	case "windows":
		return InitWindows
	default:
		return InitUnix
	}
}
