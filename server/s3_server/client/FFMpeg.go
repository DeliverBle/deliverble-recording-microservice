package client

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func ChangeFileNameMp3ToWebm(inputFileName string) (*string, error) {
	// Check if the input file exists
	if _, err := os.Stat(inputFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", inputFileName)
	}

	// Replace ".mp3" with ".webm" using a regular expression
	re := regexp.MustCompile(`\.mp3$`)
	outputFileName := re.ReplaceAllString(inputFileName, ".webm")

	// Rename the input file to the output file
	err := os.Rename(inputFileName, outputFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to rename file %s to %s: %v", inputFileName, outputFileName, err)
	}

	return &outputFileName, nil
}

func ConvertWebmBlobToMp3File(input string) error {
	cmd := exec.Command("ffmpeg", "-i", input+".webm", input+".mp3")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return fmt.Errorf("ffmpeg command failed")
	}
	fmt.Println("Result: " + out.String())

	return nil
}
