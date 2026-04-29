package parser

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Parse(pageURL string, body []byte) (string, []string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", nil, err
	}

	base, err := url.Parse(pageURL)
	title := strings.TrimSpace(doc.Find("title").First().Text())

	if err != nil {
		return title, nil, err
	}

	var links []string
	seen := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		// Skip junk schemes (this is ai generated)
		if strings.HasPrefix(href, "#") ||
			strings.HasPrefix(href, "mailto:") ||
			strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "tel:") {
			return
		}

		ref, err := url.Parse(href)
		if err != nil {
			return
		}

		absolute := base.ResolveReference(ref).String()

		if !seen[absolute] {
			seen[absolute] = true
			links = append(links, absolute)
		}
	})

	return title, links, nil
}
