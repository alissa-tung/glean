package embed

import (
	_ "embed"
	"runtime"
)

//go:embed elan-init.ps1
var InitWindows string

//go:embed elan-init.sh
var InitUnix string

func InitScriptName() string {
	switch runtime.GOOS {
	case "windows":
		return "elan-init.ps1"
	default:
		return "elan-init.sh"
	}
}

func InitScriptBytes() string {
	switch runtime.GOOS {
	case "windows":
		return InitWindows
	default:
		return InitUnix
	}
}
