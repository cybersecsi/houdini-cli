package main

import (
	"github.com/cybersecsi/houdini-cli/cmd"
	"github.com/cybersecsi/houdini-cli/pkg/houdini"
	"github.com/cybersecsi/houdini-cli/pkg/utils"
)

var version = "0.0.2"

func main() {
	utils.Banner(version)
	houdini.CheckAndCreateHoudiniDir()
	houdini.DownloadLibrary()
	houdini.LoadTools()
	cmd.Execute()
}
