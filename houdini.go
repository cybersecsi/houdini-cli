package main

import (
	"github.com/cybersecsi/houdini-cli/cmd"
	"github.com/cybersecsi/houdini-cli/pkg/houdini"
	"github.com/cybersecsi/houdini-cli/pkg/utils"
)

func main() {
	utils.Banner()
	houdini.CheckAndCreateHoudiniDir()
	houdini.DownloadToolsFile()
	houdini.LoadTools()
	cmd.Execute()
}
