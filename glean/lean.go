package glean

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func InstallLean() {
	releaseName := buildReleaseName(*version)
	resourceUrl := urlBase + "/leanprover/lean4/releases/download/v" + *version + "/" + releaseName
	log.Println("will get `" + resourceUrl + "`")

	response, err := http.Get(resourceUrl)
	if err != nil {
		panic("http.Get error: " + err.Error() + ", resourceUrl = `" + resourceUrl + "`")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	toolChainDir := filepath.Join(dotElanBaseDir, "toolchains", buildToolChainDirName(*version))
	err = os.MkdirAll(toolChainDir, 0755)
	if err != nil {
		panic(err.Error())
	}

	filePath := filepath.Join(dotElanBaseDir, "toolchains", releaseName)
	log.Println("download contents will be written to `" + filePath + "`")
	file, err := os.Create(filePath)
	if err != nil {
		panic("os.Create: " + err.Error() + ", " + "filePath = `" + filePath + "`")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("Expand-Archive", "-Force", "'"+filePath+"'", "'"+toolChainDir+"'")
	default:
		cmd = exec.Command("unzip", "-f", filePath, "-d", toolChainDir)
	}

	log.Println("exec " + cmd.String())

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	_ = os.Remove(filePath)
}

func buildReleaseName(version string) string {
	name := fmt.Sprintf("lean-%s-%s.zip", version, buildReleasePostfix())
	log.Println("buildReleaseName: `" + name + "`")
	return name
}

func buildReleasePostfix() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	default:
		switch runtime.GOARCH {
		case "aarch64":
			return runtime.GOOS + "_" + runtime.GOARCH
		default:
			return runtime.GOOS
		}
	}
}

func buildToolChainDirName(version string) string {
	return fmt.Sprintf("leanprover--lean4---%s", version)
}
