package parser

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// ExtractTitle extracts first level heading from Markdown document
func ExtractTitle(data []byte) string {
	var title string

	ast.WalkFunc(parse(data), func(node ast.Node, entering bool) ast.WalkStatus {
		if heading, ok := node.(*ast.Heading); ok {
			if entering && heading.Level == 1 {
				title = joinHeadingLiterals(heading)
				return ast.Terminate
			}
		}
		return ast.GoToNext
	})

	return title
}

func joinHeadingLiterals(heading *ast.Heading) string {
	var title string
	for _, c := range heading.GetChildren() {
		leaf := c.AsLeaf()
		if leaf == nil {
			continue
		}
		title = fmt.Sprintf("%s%s", title, leaf.Literal)
	}

	return title
}

// ExtractLinks extracts unique links from Markdown document
func ExtractLinks(data []byte) []string {
	var links []string
	linksMap := make(map[string]struct{})
	d := parse(data)

	ast.WalkFunc(d, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node.(type) {
		case *ast.HTMLBlock:
			ll, err := extractUniqHTMLLinks(string(node.AsLeaf().Literal))
			if err != nil {
				break
			}
			for l := range ll {
				linksMap[l] = struct{}{}
			}
		case *ast.HTMLSpan:
			ll, err := extractUniqHTMLLinks(string(node.AsLeaf().Literal))
			if err != nil {
				break
			}
			for l := range ll {
				linksMap[l] = struct{}{}
			}
		case *ast.Link:
			linksMap[string(node.(*ast.Link).Destination)] = struct{}{}
		case *ast.Image:
			linksMap[string(node.(*ast.Image).Destination)] = struct{}{}
		}

		return ast.GoToNext
	})

	for l := range linksMap {
		links = append(links, l)
	}

	return links
}

// parse creating a new parser and make parsing.
// Each time parser is new because parsing change its state.
func parse(input []byte) ast.Node {
	return parser.New().Parse(input)
}
