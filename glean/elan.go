package glean

import (
	"fmt"
	"github.com/alissa-tung/glean/embed"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func InstallElan() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	scriptName := "glean.tmp." + embed.InitScriptName()
	scriptPath := filepath.Join(cwd, scriptName)
	log.Println("write init script to `" + scriptPath + "`")
	file, err := os.Create(scriptPath)
	if err != nil {
		panic("Failed to create `" + scriptName + "`, " + err.Error())
	}

	func() {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		_, err = file.Write(embed.InitScriptBytes())
		if err != nil {
			panic(err)
		}
	}()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-f", scriptPath)
	default:
		cmd = exec.Command("/bin/sh", scriptPath, "-y", "--default-toolchain", "none")

		zshText := `
########
# begin @generated by glean
PATH="$HOME/.elan/bin:$PATH"
# end @generated by glean
########
`
		zprofilePath := path.Join(os.Getenv("HOME"), ".zprofile")
		log.Println("Adding glean to zsh PATH")
		file, err := os.OpenFile(zprofilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("append to zprofile error: " + err.Error())
		}
		_, err = file.WriteString(zshText)
		if err != nil {
			log.Fatalf("append to zprofile error: " + err.Error())
		}
	}

	log.Println("exec `" + cmd.String() + "`")

	o, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(o))

	_ = os.Remove(scriptPath)
}
