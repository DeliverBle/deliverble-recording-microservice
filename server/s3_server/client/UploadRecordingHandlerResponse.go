package client

type UploadRecordingHandlerResponse struct {
	Code int    `json:"code"`
	Url  string `json:"url"`
	Key  string `json:"key"`
}
