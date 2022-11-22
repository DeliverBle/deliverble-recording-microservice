package test

import (
	"bytes"
	"deliverble-recording-msa/preprocess"
	"deliverble-recording-msa/server/s3_server/client"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadRecordingHandler(t *testing.T) {
	/*
		Run `go build deliverble-recrording-msa/server/s3_server` before running this test.
	*/

	// given : preparing form file (mp3)
	assertions := assert.New(t)

	path := "./test.mp3"
	file, err := os.Open(path)
	assertions.NoError(err)
	defer file.Close()

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	multi, err := writer.CreateFormFile("file", filepath.Base(path))
	assertions.NoError(err)

	_, err = io.Copy(multi, file)
	assertions.NoError(err)

	err = writer.Close()
	assertions.NoError(err)

	// when
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", buf)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	// create new c echo.Context for testing
	con := echo.New().NewContext(req, res)

	// then
	err = preprocess.UploadRecordingHandler(con)
	assertions.NoError(err)

	body, err := ioutil.ReadAll(res.Body)
	assertions.NoError(err)

	var response client.UploadRecordingHandlerResponse
	err = json.Unmarshal(body, &response)

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, response.Url, "https://deliverable-recording.s3.ap-northeast-2.amazonaws.com")
	assert.Contains(t, response.Url, ".mp3")
}
