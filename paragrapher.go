package tagfunctions

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// paragrapher -- creates paragraphs tags based on \n\n

func Paragrapher(n *html.Node) error {

	blocks := Selector(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "blockquote"
	})
	for _, block := range blocks {
		p := &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.P,
			Data:     "p",
		}
		bq := &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Blockquote,
			Data:     "blockquote",
		}
		current := block.FirstChild
		for current != nil {
			// TODO: if it's a block element, close of paragraph tag if any
			// and add it.
			if current.Type == html.ElementNode {
				next := current.NextSibling
				block.RemoveChild(current)
				bq.AppendChild(current)
				current = next
				continue
			}
			if current.Type == html.TextNode {
				ps := strings.Split(current.Data, "\n\n")
				p.AppendChild(&html.Node{
					Type: html.TextNode,
					Data: ps[0],
				})
				for _, txt := range ps[1:] {
					bq.AppendChild(p)

					p = &html.Node{
						Type:     html.ElementNode,
						DataAtom: atom.P,
						Data:     "p",
					}
					p.AppendChild(&html.Node{
						Type: html.TextNode,
						Data: txt,
					})
				}

				next := current.NextSibling
				block.RemoveChild(current)
				current = next
				continue
			}
		}
		if p.FirstChild != nil {
			bq.AppendChild(p)
		}
		// at this point we have a new blockquote element all
		// filled out
		parent := block.Parent
		parent.InsertBefore(bq, block)
		parent.RemoveChild(block)
	}

	return nil
}
