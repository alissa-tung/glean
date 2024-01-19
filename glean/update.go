package glean

import (
	"fmt"
	"os"
	"os/exec"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"runtime"
	"log"
	"path/filepath"
)

func buildGleanReleaseName() string {
	var arch string
	if runtime.GOARCH == "amd64" {
		arch = "x86_64"
	} else{
		arch = "arm64"
	}
	
	switch runtime.GOOS {
	case "windows":
		name := fmt.Sprintf("glean_%s_%s.zip", "Windows", arch)
		return name
	default:
		name := fmt.Sprintf("glean_%s_%s.tar.gz", strings.Title(runtime.GOOS), arch)
		return name
	}
}

func getLatestVersion() string {
	
	response, err := http.Get(urlBase + "/glean/releases/download/?mirror_intel_list")
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	bodyText := string(body)
	index := strings.Index(bodyText, "v0")
	var latestVersion string
	for index != -1 {
		latestVersion = bodyText[index:index+6]
		fmt.Println("Found latest version:", latestVersion)
		break
	}
	return latestVersion
}


func CheckUpdate() {
	fmt.Println("Checking for updates...")
	latestVersion := getLatestVersion()
	if latestVersion == gleanVersion {
		fmt.Println("Already up to date")
		return
	}
	fmt.Println("New version available:", latestVersion)

	releaseName := buildGleanReleaseName()
	response, err := http.Get(urlBase + "/glean/releases/download/" + latestVersion + "/" + releaseName)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	filePath := filepath.Join(dotElanBaseDir, "bin", releaseName)
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
		panic(err.Error())
	}
	file.Close()
	var cmd *exec.Cmd
	gleantmpPath := filepath.Join(dotElanBaseDir, "bin", "glean.new")
	err = os.MkdirAll(gleantmpPath, 0755)
	if err != nil {
		panic(err.Error())
	}

	cmd = exec.Command("tar", "-xvf", filePath, "-C", gleantmpPath)

	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}
	err = os.Remove(filePath)
	if err != nil {
		panic(err.Error())
	}
	switch runtime.GOOS {
	case "windows":
		fmt.Printf("Please run the command `cp %s\\glean.exe %s`", gleantmpPath, dotElanBaseDir+"\\bin")
	default:
		fmt.Printf("Please run the command `cp %s/glean %s`", gleantmpPath, dotElanBaseDir+"/bin")
	}
	os.Exit(0)
}