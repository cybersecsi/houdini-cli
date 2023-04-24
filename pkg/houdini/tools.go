package houdini

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cybersecsi/houdini-cli/pkg/types"
	"github.com/cybersecsi/houdini-cli/pkg/utils"
)

func LoadTools() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	houdiniDir := filepath.Join(homeDir, ".houdini")
	libraryDirPath := filepath.Join(houdiniDir, "library")

	tools, err := ioutil.ReadDir(libraryDirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, toolFolder := range tools {
		toolJsonConfig := filepath.Join(houdiniDir, "library", toolFolder.Name(), "config.json")
		// Read the contents of the JSON file
		jsonFile, err := ioutil.ReadFile(toolJsonConfig)
		if err != nil {
			log.Fatal(err)
		}
		var tool types.Tool
		// Parse the JSON data and push into the array of Tool objects
		err = json.Unmarshal(jsonFile, &tool)
		if err != nil {
			log.Fatal(err)
		}
		utils.Tools = append(utils.Tools, tool)
	}
}

func GetToolNames() []string {
	names := []string{}
	for _, tool := range utils.Tools {
		toolName := fmt.Sprintf("%s/%s", tool.Organization, tool.Name)
		names = append(names, toolName)
	}
	return names
}

func ListTools() {
	for _, tool := range utils.Tools {
		fmt.Printf("%s/%s\n", tool.Organization, tool.Name)
	}
}

func FindTool(fullName string) (*types.Tool, error) {
	res := strings.Split(fullName, "/")
	org := res[0]
	name := res[1]
	for _, tool := range utils.Tools {
		if tool.Name == name && tool.Organization == org {
			return &tool, nil
		}
	}
	return nil, errors.New("Unable to find the tool")
}
