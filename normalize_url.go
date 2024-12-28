package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawUrl string) (string, error) {
	var normalizedURL string
	url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	urlPath := strings.TrimSuffix(strings.TrimPrefix(url.Path, "/"), "/")
	normalizedURL = fmt.Sprintf("%s/%s", url.Host, urlPath)
	return normalizedURL, nil
}
