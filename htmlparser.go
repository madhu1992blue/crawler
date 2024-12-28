package main

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getLinksFromNode(node *html.Node) []string {
	var links []string
	if node.Type != html.ElementNode {
		return links
	}
	if node.DataAtom != atom.A {
		for child := range node.ChildNodes() {
			links = append(links, getLinksFromNode(child)...)
		}
		return links
	}
	nodeAttributes := node.Attr
	for _, attr := range nodeAttributes {
		if attr.Key == "href" {
			parsedURL, err := url.Parse(attr.Val)
			if err != nil {
				log.Printf("Couldn't parse URL: %s\n", parsedURL)
				continue
			}
			links = append(links, attr.Val)
		}
	}
	return links

}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var urls []string
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	htmlRoot, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}
	childNodes := htmlRoot.ChildNodes()
	for node := range childNodes {
		log.Println(node.DataAtom.String())
		urls = append(urls, getLinksFromNode(node)...)
	}
	for i, rawUrl := range urls {
		url, err := url.Parse(rawUrl)
		if err != nil {
			return nil, err
		}

		urls[i] = baseURL.ResolveReference(url).String()

	}
	return urls, nil
}
