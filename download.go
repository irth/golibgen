package libgen

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getDownloadLink(mirrors []string) (string, error) {
	var mirrorURL *url.URL = nil
	for _, mirror := range mirrors {
		u, err := url.Parse(mirror)
		if err != nil {
			continue
		}
		if u.Host == "library.lol" {
			mirrorURL = u
		}
	}
	if mirrorURL == nil {
		return "", fmt.Errorf("no supported mirrors found")
	}

	switch mirrorURL.Host {
	case "library.lol":
		return getDownloadLinkLibraryLol(mirrorURL)
	default:
		return "", fmt.Errorf("Unsupported mirror: %s", mirrorURL.Host)
	}
}

func getDownloadLinkLibraryLol(mirrorURL *url.URL) (string, error) {
	res, err := http.Get(mirrorURL.String())
	if err != nil {
		return "", fmt.Errorf("couldn't get the mirror page: %w", err)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't parse the mirror page: %w", err)
	}

	body := findBody(doc)
	if body == nil {
		return "", fmt.Errorf("parse error: cannot find body")
	}

	link := findElement(body, func(el *html.Node) bool {
		return el.Type == html.ElementNode && el.Data == "a" && strings.TrimSpace(getText(el)) == "Cloudflare"
	})
	if link == nil {
		return "", fmt.Errorf("couldn't find supported download link")
	}
	linkStr, ok := getAttr(link, "href")
	if !ok {
		return "", fmt.Errorf("download link has no href attribute")
	}
	return linkStr, nil
}
