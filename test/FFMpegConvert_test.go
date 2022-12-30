package test

import (
	ffw "deliverble-recording-msa/server/s3_server/client"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func copyFile(src, dst string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func TestFFMpegConvert(t *testing.T) {
	// given: file name
	Mp3Extension := ".mp3"
	WebMExtension := ".webm"
	SampleSourceName := "test" // source mp3
	SrcName := "test2"         // fake mp3 (copied from source mp3 for testing)
	DstName := "output2"       // real mp3

	// given: copy test.mp3 to test2.mp3
	errCopy := copyFile(SampleSourceName+Mp3Extension, SrcName+Mp3Extension)
	if errCopy != nil {
		fmt.Println(errCopy)
	}
	assert.NoError(t, errCopy)

	assertions := assert.New(t)
	if _, err := os.Stat(DstName + Mp3Extension); err == nil {
		err = os.Remove(DstName + Mp3Extension)
		assertions.NoError(err)
	}

	// when: ChangeFileNameMp3ToWebm
	outputFileNameWithWebM, errChange := ffw.ChangeFileNameMp3ToWebm(SrcName + Mp3Extension)
	if errChange != nil {
		fmt.Println(errChange)
		assertions.Fail("failed to change file name : ", errChange.Error())
	}

	// when: ConvertWebmBlobToMp3File
	errConvert := ffw.ConvertWebmBlobToMp3File(SrcName)
	if errConvert != nil {
		fmt.Println(errConvert)
		assertions.Fail("ffmpeg command failed : ", errConvert.Error())
	}

	// then
	assertions.NoError(errChange)
	assertions.NoError(errConvert)
	assertions.Equal(SrcName+WebMExtension, *outputFileNameWithWebM)
	assertions.FileExists(SrcName + Mp3Extension)
}
