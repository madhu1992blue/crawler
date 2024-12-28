package main

import (
	"errors"
	"io"
	"mime"
	"net/http"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", errors.New("client error Getting data")
	}
	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return "", errors.New("Content-Type malformed")
	}
	if mediaType != "text/html" {
		return "", errors.New("response not html")
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBodyBytes), nil

}
