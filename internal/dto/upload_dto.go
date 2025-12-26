package dto

// ------ body ------
type PresignUploadBody struct {
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type CommitUploadBody struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Tags     []struct {
		KeyID   int32 `json:"keyID"`
		ValueID int32 `json:"valueID"`
	} `json:"tags"`
}

type ReplaceTagsBody struct {
	Tags []TagPair `json:"tags"`
}

// ------ request ------
type PresignUploadRequest struct {
	TeamPathParams
	PresignUploadBody
}

type CommitUploadRequest struct {
	TeamPathParams
	CommitUploadBody
}

type ReplaceTagsRequest struct {
	UploadRefPathParams
	ReplaceTagsBody
}

// ------ response ------
type PresignUploadResponse struct {
	Url    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
