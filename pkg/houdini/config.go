package houdini

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cybersecsi/houdini-cli/pkg/utils"
)

func CheckAndCreateHoudiniDir() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	houdiniDir := filepath.Join(homeDir, ".houdini")

	if _, err := os.Stat(houdiniDir); os.IsNotExist(err) {
		utils.Info(fmt.Sprintf("Creating %s/.houdini directory...", homeDir))
		if err := os.Mkdir(houdiniDir, 0700); err != nil {
			utils.Error(fmt.Sprintf("Failed to create %s/.houdini directory", homeDir))
			os.Exit(1)
		}
	}
}

func DownloadLibrary() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	houdiniDir := filepath.Join(homeDir, ".houdini")
	libraryFilePath := filepath.Join(houdiniDir, "houdini-library.tar.gz")
	repoURL := "https://api.github.com/repos/cybersecsi/houdini/releases/latest"

	// Make a GET request to the release endpoint
	resp, err := http.Get(repoURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Parse the JSON response to get the download URL for the latest release asset
	// Replace "asset_name" with the name of the asset you want to download
	// Alternatively, you can loop through the assets and download all of them
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return
	}
	downloadURL := result["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"].(string)

	// Download the asset
	err = utils.DownloadFile(downloadURL, libraryFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.UnpackTarGz(libraryFilePath, houdiniDir)

	err = os.Remove(libraryFilePath)
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}
}

func UpdateToolsFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	houdiniDir := filepath.Join(homeDir, ".houdini")
	libraryFolderPath := filepath.Join(houdiniDir, "library")

	if _, err := os.Stat(libraryFolderPath); err == nil {
		err := os.RemoveAll(libraryFolderPath)
		if err != nil {
			fmt.Println("Error deleting folder:", err)
			return
		}
	}
	DownloadLibrary()
}
