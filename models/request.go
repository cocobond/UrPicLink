package models

type UploadImageRequest struct {
	Base64 string `json:"base64" binding:"required"`
	Mime   string `json:"mime" binding:"required"`
}

type UploadFileRequest struct {
	File interface{} `json:"file" binding:"required"`
}

type GithubAPIRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

type GithubAPIResponse struct {
	Content struct {
		DownloadURL string `json:"download_url"`
	} `json:"content"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}
