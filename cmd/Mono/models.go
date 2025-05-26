package main

type ExtractResponse struct {
	Title   string
	Content string
}

type ConversionRequest struct {
	Html  string `json:"html"`
	Title string `json:"title"`
}

type UploadOptions struct {
	Token string
	Path  string
	Data  []byte
}
