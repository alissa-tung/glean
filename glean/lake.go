package glean

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type gitRepoSourceMirrorMapping struct {
	sourceUrl string
	mirrorUrl string
}

func buildSjtuUrl(pkgName string) string {
	return fmt.Sprintf("https://mirror.sjtu.edu.cn/git/lean4-packages/%s.git", pkgName)
}

func getRepoName(url string) string {
	xs := strings.Split(url, "/")
	xs = strings.Split(xs[len(xs)-1], ".")
	return xs[0]
}

func buildSourceToSjtu(url string) gitRepoSourceMirrorMapping {
	return gitRepoSourceMirrorMapping{
		url,
		buildSjtuUrl(getRepoName(url)),
	}
}

var (
	repos = [...]string{
		"https://github.com/JLimperg/aesop.git",
		"https://github.com/leanprover/doc-gen4.git",
		"https://github.com/leanprover/lean4-cli.git",
		"https://github.com/avigad/mathematics_in_lean_source.git",
		"https://github.com/leanprover-community/mathlib4.git",
		"https://github.com/leanprover-community/ProofWidgets4.git",
		"https://github.com/gebner/quote4.git",
		"https://github.com/leanprover/std4.git",
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
	url  string `json:"url"`
	rev  string `json:"rev"`
	name string `json:"name"`
}

type lakeGitPackage struct {
	git lakePackage `json:"git"`
}

type lakeManifest struct {
	packages map[string]interface{} `json:"packages"`
}

func readAndParse(url string) lakeManifest {
	file, err := os.ReadFile(url)
	if err != nil {
		panic("reading `" + url + "`, " + err.Error())
	}

	var obj lakeManifest
	if err = json.Unmarshal(file, &obj); err != nil {
		panic("error parsing lake manifest `" + url + "`, " + err.Error())
	}
	return obj
}

func LakeSyncPackages() {
	obj := readAndParse(*LakeManifestPath)
	fmt.Println(obj)

	var reposToClone []lakePackage

	//for _, v := range obj.packages {
	//	if mirrorUrl := findMirror(v.git.url); mirrorUrl != "" {
	//		v.git.url = mirrorUrl
	//		reposToClone = append(reposToClone, v.git)
	//	}
	//}

	for _, v := range reposToClone {
		log.Printf("repo to clone: %+v\n", v)
	}
}

func findMirror(url string) string {
	for _, v := range mirrorRepos {
		if v.sourceUrl == url {
			return v.mirrorUrl
		}
	}
	return ""
}
