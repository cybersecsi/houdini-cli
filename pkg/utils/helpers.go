package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func CompileRunCommand(runCommand string) string {
	// Sample string with tags

	// Define a regular expression to match tags
	re := regexp.MustCompile("<(\\w+)>")

	// Find all matches in the string
	matches := re.FindAllStringSubmatch(runCommand, -1)

	// Create a map to store user input for each tag
	inputs := make(map[string]string)

	// Ask the user to set a value for each tag
	reader := bufio.NewReader(os.Stdin)
	for _, match := range matches {
		tag := match[1]
		fmt.Printf("Enter a value for '%s%s%s': ", boldStart, tag, boldEnd)
		value, _ := reader.ReadString('\n')
		inputs[tag] = strings.TrimSpace(value)
	}

	// Replace tags with user input
	for tag, value := range inputs {
		runCommand = strings.ReplaceAll(runCommand, fmt.Sprintf("<%s>", tag), value)
	}

	// Return the final string with replaced values
	return runCommand
}

func ConfirmAction(actionTxt string) bool {
	actionConfirmed := true
	prompt := &survey.Confirm{
		Message: actionTxt,
		Default: true,
	}
	err := survey.AskOne(prompt, &actionConfirmed)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return actionConfirmed
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnpackTarGz(tarGzFile string, destDir string) error {
	// Open the tar.gz file
	f, err := os.Open(tarGzFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a gzip reader for the file
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	// Create a tar reader for the gzip reader
	tr := tar.NewReader(gz)

	// Iterate over the files in the tar archive
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return err
		}

		// Create the destination file
		destFile := destDir + "/" + hdr.Name
		fi := hdr.FileInfo()
		if fi.IsDir() {
			err = os.MkdirAll(destFile, fi.Mode())
			if err != nil {
				return err
			}
			continue
		}
		f, err := os.OpenFile(destFile, os.O_CREATE|os.O_RDWR, fi.Mode())
		if err != nil {
			return err
		}

		// Copy the contents of the file from the tar archive to the destination file
		_, err = io.Copy(f, tr)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}
