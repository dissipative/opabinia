package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// extractUniqHTMLLinks extracts links from HTML string and returns a map of unique links
func extractUniqHTMLLinks(html string) (map[string]struct{}, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	links := make(map[string]struct{})

	extractBySelector(doc, "a", "href", links)
	extractBySelector(doc, "img", "src", links)
	extractBySelector(doc, "audio", "src", links)
	extractBySelector(doc, "video", "src", links)

	return links, nil
}

// extractBySelector extracts links from HTML document by selector and attribute name
func extractBySelector(doc *goquery.Document, selector, attrName string, links map[string]struct{}) {
	_ = doc.Find(selector).Each(func(num int, sel *goquery.Selection) {
		src := sel.AttrOr(attrName, "")
		if src != "" {
			links[src] = struct{}{}
		}
	})
}
