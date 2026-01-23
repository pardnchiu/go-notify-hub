package email

type Email struct {
	Text         string `json:"text,omitempty"`
	Image        string `json:"image,omitempty"`
	ImagePreview string `json:"image_preview,omitempty"`
}
