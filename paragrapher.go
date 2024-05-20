package tagfunctions

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// paragrapher -- creates paragraphs tags based on \n\n
//
// easy case: one text child
// <p>line1\n\nline2</p>
// insert <p>line1</p><p>line2</p> before orig
// delete orig
//
// <p>line1\n\nfoo<b>line2</b>junk</p>
// <p>line1</p><p>foo<b>line2</b>junk</p>
//
// $p{$b{bold} line1 \n\n line2} --> $p{$b{bold} line1} $p{line2}
// $p{$b{bold} line1 \n\n $b{line2}} -> $p{$b{boldi} line1} $p{$b{line2}}
// $p{$b{bold} line1 \n\n line2 $b{bold}} -->$p{$b{bold} line1} $p{line2 $b{bold}}
//
// Since this is a paragraph, whitespace is trimmed from front and start of
// paragraph.  This also means repeated runs of "\n" only make on paragraph
//
//	split.
func isBlock(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	switch n.Data {
	case "p", "table", "blockquote", "pre":
		return true

	}
	return false
}
func Paragrapher(n *html.Node) error {
	if err := ParagrapherTag(n, "p"); err != nil {
		return err
	}
	return ParagrapherTag(n, "blockquote")
}
func ParagrapherTag(n *html.Node, tagName string) error {
	tagAtom := atom.Lookup([]byte(tagName))
	blocks := Selector(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagName
	})
	for _, block := range blocks {

		p := &html.Node{
			Type:     html.ElementNode,
			DataAtom: tagAtom,
			Data:     tagName,
		}
		current := block.FirstChild
		for current != nil {
			// generate case of <p> having block-level children (e.e. <p> outer <p> inner </p></p>
			if isBlock(current) {
				if p.FirstChild != nil {
					block.Parent.InsertBefore(p, block)
					p = &html.Node{
						Type:     html.ElementNode,
						DataAtom: tagAtom,
						Data:     tagName,
					}
				}
				next := current.NextSibling
				block.RemoveChild(current)
				block.Parent.InsertBefore(current, block)
				current = next
				continue
			}
			// if not text, remove and copy to next container
			if current.Type != html.TextNode {
				next := current.NextSibling
				block.RemoveChild(current)
				p.AppendChild(current)
				current = next
				continue
			}
			// everything is a TextNode now
			// special case: empty text block.. skip and remove empty
			if current.Data == "" {
				next := current.NextSibling
				block.RemoveChild(current)
				current = next
				continue
			}
			idx := strings.Index(current.Data, "\n\n")

			// no "\n\n", so copy text node over to container
			if idx == -1 {
				next := current.NextSibling
				block.RemoveChild(current)
				p.AppendChild(current)
				current = next
				continue
			}

			// we are textnode and has a "\n\n"
			p.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: strings.TrimSpace(current.Data[:idx]),
			})
			block.Parent.InsertBefore(p, block)
			p = &html.Node{
				Type:     html.ElementNode,
				DataAtom: tagAtom,
				Data:     tagName,
			}
			current.Data = strings.TrimSpace(current.Data[idx+2:])
		}

		if p.FirstChild != nil {
			block.Parent.InsertBefore(p, block)
		}
		block.Parent.RemoveChild(block)

	}
	return nil
}
