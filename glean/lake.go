package glean

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type gitRepoSourceMirrorMapping struct {
	sourceUrl string
	mirrorUrl string
}

func buildSjtuUrl(pkgName string) string {
	return fmt.Sprintf("https://mirror.sjtu.edu.cn/git/lean4-packages/%s", pkgName)
}

func getRepoName(url string) string {
	xs := strings.Split(url, "/")
	x := xs[len(xs)-1]
	return x
}

func buildSourceToSjtu(url string) gitRepoSourceMirrorMapping {
	return gitRepoSourceMirrorMapping{
		url,
		buildSjtuUrl(getRepoName(url)),
	}
}

var (
	repos = [...]string{
		"https://github.com/JLimperg/aesop",
		"https://github.com/leanprover-community/aesop",
		"https://github.com/leanprover/doc-gen4",
		"https://github.com/leanprover/lean4-cli",
		"https://github.com/mhuisi/lean4-cli",
		"https://github.com/avigad/mathematics_in_lean_source",
		"https://github.com/leanprover-community/mathlib4",
		"https://github.com/leanprover-community/ProofWidgets4",
		"https://github.com/EdAyers/ProofWidgets4",
		"https://github.com/gebner/quote4",
		"https://github.com/leanprover-community/quote4",
		"https://github.com/leanprover/std4",
		"https://github.com/leanprover-community/import-graph",
		"https://github.com/leanprover-community/batteries"
	}

	mirrorRepos = func() []gitRepoSourceMirrorMapping {
		var ret []gitRepoSourceMirrorMapping
		for _, v := range repos {
			ret = append(ret, buildSourceToSjtu(v))
		}
		return ret
	}()
)

func projectDir() string {
	return filepath.Dir(*LakeManifestPath)
}

type lakePackage struct {
	Url  string `json:"url"`
	Rev  string `json:"rev"`
	Name string `json:"name"`
}

type lakeManifest struct {
	Version     int           `json:"version"`
	PackagesDir string        `json:"packagesDir"`
	Packages    []lakePackage `json:"packages"`
}

func readAndParse(url string) lakeManifest {
	file, err := os.ReadFile(url)
	if err != nil {
		fmt.Println("reading `" + url + "`, " + err.Error())
		os.Exit(0)
	}

	var obj lakeManifest
	if err = json.Unmarshal(file, &obj); err != nil {
		panic("error parsing lake manifest `" + url + "`, " + err.Error())
	}
	return obj
}

func LakeSyncPackages() {
	obj := readAndParse(*LakeManifestPath)

	var reposToClone []lakePackage

	if len(obj.Packages) == 0 {
		panic("empty packages in manifest json")
	}
	for _, v := range obj.Packages {
		mirrorUrl := findMirror(v.Url)
		if mirrorUrl != "" {
			v.Url = mirrorUrl
			reposToClone = append(reposToClone, v)
		} else {
			log.Println("failed to find mirror for: `" + v.Url + "`")
		}
	}

	for _, v := range reposToClone {
		log.Printf("repo to clone: %v\n", v.Name)

		target := filepath.Join(projectDir(), obj.PackagesDir, v.Name)
		if err := os.RemoveAll(target); err != nil {
			log.Println("Failed to remove `" + target + "`, " + err.Error())
		}
		cmd := exec.Command("git", "clone", v.Url, target)
		if err := cmd.Run(); err != nil {
			panic("Failed to clone `" + cmd.String() + "`, " + err.Error())
		}
		cmd = exec.Command("git", "-C", target, "checkout", v.Rev)
		if err := cmd.Run(); err != nil {
			panic("Failed to checkout `" + cmd.String() + "`, " + err.Error())
		}

	}
}

func findMirror(url string) string {
	for _, v := range mirrorRepos {

		if v.sourceUrl == url || v.sourceUrl+".git" == url {
			return v.mirrorUrl
		}
	}
	return ""
}
