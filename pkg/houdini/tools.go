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
	"sync"

	"github.com/cybersecsi/houdini-cli/pkg/types"
	"github.com/cybersecsi/houdini-cli/pkg/utils"
)

func loadTool(wg *sync.WaitGroup, toolName string) {
	defer wg.Done()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	houdiniDir := filepath.Join(homeDir, ".houdini")
	toolJsonConfig := filepath.Join(houdiniDir, "library", toolName, "config.json")
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

	concurrency := 10
	var wg sync.WaitGroup
	wg.Add(len(tools))

	jobs := make(chan string)
	for i := 0; i < concurrency; i++ {
		go func() {
			for element := range jobs {
				loadTool(&wg, element)
			}
		}()
	}

	// Add jobs to the channel
	for _, toolFolder := range tools {
		jobs <- toolFolder.Name()
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

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
