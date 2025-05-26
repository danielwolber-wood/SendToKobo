package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

type UploadOptions struct {
	Token string
	Path  string
	Data  []byte
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("POST /v1/api/upload", handleUpload)
	http.ListenAndServe(":80", r)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		msg := fmt.Sprintf("failed to parse data: %v\n", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	token := r.FormValue("token")
	filename := r.FormValue("filename")
	path := r.FormValue("filepath")

	file, header, err := r.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("failed to get file data: %v\n", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer file.Close()
	if filename == "" {
		filename = header.Filename
	}

	fullPath := filepath.Join(path, filename)

	data, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("failed to read file data: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	opts := UploadOptions{
		Token: token,
		Path:  fullPath,
		Data:  data,
	}

	err = Upload(opts)
	if err != nil {
		msg := fmt.Sprintf("failed to upload file: %v\n", err)
		http.Error(w, msg, http.StatusInternalServerError) // TODO wrong error type
	}
}

func Upload(opts UploadOptions) error {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", bytes.NewReader(opts.Data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf("{\"path\": \"%s\"}", opts.Path))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", opts.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status: %v\n", resp.Status)
	fmt.Printf("body: %v\n", body)
	return nil
}
