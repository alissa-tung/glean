package main

import (
	"flag"
	"fmt"
	"github.com/alissa-tung/glean/glean"
)

func main() {
	glean.InitFlags()
	glean.GetLatestVersion()
	fmt.Println("Please refer to `https://mirror.sjtu.edu.cn/elan/?mirror_intel_list` for available versions")
	if *glean.Update {
		glean.CheckUpdate()
	}
	switch *glean.Command {
	case "elan":
		glean.InstallElan()

	case "lean":
		glean.InstallLean()

	default:
		if *glean.LakeManifestPath != "" {
			glean.LakeSyncPackages()
		} else {
			fmt.Println("unknown command")
			flag.Usage()
		}
	}
}
