package request

type UploadPresignedURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type UploadPresignedURLsRequest struct {
	Files []UploadPresignedURLRequest `json:"files" binding:"required,min=1,dive"`
}

type ViewPresignedURLsRequest struct {
	Keys []string `json:"keys" binding:"required,min=1,dive"`
}
