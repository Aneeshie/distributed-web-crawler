package parser

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Parse(body []byte) (string, []string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", nil, err
	}

	title := strings.TrimSpace(doc.Find("title").First().Text())

	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.TrimSpace(href) != "" {
			links = append(links, href)
		}
	})

	return title, links, nil

}
