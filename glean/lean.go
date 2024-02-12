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
	"strings"
)

func InstallLean() {
	releaseName := buildReleaseName(*version)
	var resourceUrl string
	if strings.Contains(*version, "nightly") {
		nightlyIndex := strings.Index(*version, "nightly")
		nightlyVersion := *version
		resourceUrl = urlBase + "/leanprover/lean4_nightly/releases/download/" + nightlyVersion[nightlyIndex:] + "/" + releaseName
	} else {
		resourceUrl = urlBase + "/leanprover/lean4/releases/download/v" + *version + "/" + releaseName
	}
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

	finalToolChainDir := filepath.Join(dotElanBaseDir, "toolchains", buildToolChainDirName(*version))
	err = os.MkdirAll(finalToolChainDir, 0755)
	if err != nil {
		panic(err.Error())
	}

	tmpToolChainDir := filepath.Join(dotElanBaseDir, "toolchains", "tmp")
	err = os.MkdirAll(tmpToolChainDir, 0755)
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
	file.Close()

	var cmd *exec.Cmd
	_ = os.RemoveAll(finalToolChainDir)
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tar", "-xvf", filePath, "-C", tmpToolChainDir)
	default:
		cmd = exec.Command("unzip", filePath, "-d", tmpToolChainDir)
	}

	log.Println("exec `" + cmd.String() + "`")

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	oldPath := filepath.Join(
		tmpToolChainDir,
		fmt.Sprintf("lean-%s-%s", *version, buildReleasePostfix()),
	)
	if err := os.Rename(oldPath, finalToolChainDir); err != nil {
		panic("Can not rename `" + oldPath + "` to `" + finalToolChainDir + "`")
	}

	if err != nil {
		return
	}

	_ = os.Remove(filePath)
	_ = os.RemoveAll(tmpToolChainDir)
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
	if strings.Contains(version, "nightly") {
		nightlyVersion := version[strings.Index(version, "nightly"):]
		return fmt.Sprintf("leanprover--lean4---%s", nightlyVersion)
	} else {
		return fmt.Sprintf("leanprover--lean4---v%s", version)
	}
}
