package Entities

type UploadLog struct {
	Id             int    `json:"id"`
	Uploader       int    `json:"uploader"`
	Path           string `json:"path"`
	Finished       bool   `json:"finished"`
	Received_bytes uint64 `json:"received_bytes"`
	Size           uint64 `json:"size"`
	ClientFilePath string `json:"client_file_path"`
}
