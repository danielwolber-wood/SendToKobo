package main

import (
	"bytes"
	"fmt"
	"github.com/go-shiori/go-readability"
	"html"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func ConvertStringWithPandoc(content, title, fromFormat, toFormat string) ([]byte, error) {
	cmd := exec.Command("pandoc", "-f", fromFormat, "-t", toFormat, "--title", title, "--quiet")

	cmd.Stdin = strings.NewReader(content)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func sanitizeHTML(content string) string {
	dangerous := []string{
		"<script", "</script>",
		"<iframe", "</iframe>",
		"javascript:",
		"data:text/html",
		"vbscript:",
	}

	cleaned := content
	for _, danger := range dangerous {
		cleaned = strings.ReplaceAll(strings.ToLower(cleaned), danger, "")
	}

	return html.UnescapeString(cleaned)
}

func Extract(r io.Reader) (ExtractResponse, error) {
	urlObj, err := url.Parse("example.com")
	if err != nil {
		return ExtractResponse{}, err
	}
	article, err := readability.FromReader(r, urlObj)
	if err != nil {
		return ExtractResponse{}, err
	}
	return ExtractResponse{Title: article.Title, Content: article.Content}, nil
}

func GenerateHTML(title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body>
    %s
</body>
</html>`, title, body)
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

func sanitizeString(input string) string {
	replacer := strings.NewReplacer(
		"?", "_",
		"\"", "_",
		"*", "_",
		"\\", "_",
		"|", "_",
		"/", "_",
	)
	return replacer.Replace(input)
}
