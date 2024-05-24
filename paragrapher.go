package tagfunctions

import (
	"strings"

	"golang.org/x/net/html"
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
	case "p", "table", "blockquote", "pre", "div":
		return true

	}
	return false
}

type Paragrapher struct {
	Tags   []string // elements to split on "\n\n" to generate new blocks
	Blocks []string // elements that are considered to be blocks
}

func (p *Paragrapher) Execute(n *html.Node) error {

	//
	// Set defaults
	//
	if len(p.Tags) == 0 {
		p.Tags = []string{"p"}
	}
	if len(p.Blocks) == 0 {
		p.Blocks = []string{"p", "blockquote"}
	}

	for _, tag := range p.Tags {
		if err := p.executeTag(n, tag); err != nil {
			return err
		}
	}

	return nil
}

func (p *Paragrapher) executeTag(n *html.Node, tagName string) error {
	blocks := Selector(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagName
	})
	for _, block := range blocks {
		p := NewElement(tagName)
		current := block.FirstChild
		for current != nil {
			// generate case of <p> having block-level children (e.e. <p> outer <p> inner </p></p>
			if isBlock(current) {
				if p.FirstChild != nil {
					block.Parent.InsertBefore(p, block)
					p = NewElement(tagName)
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

			copyText := strings.TrimSpace(current.Data[:idx])
			current.Data = strings.TrimSpace(current.Data[idx+2:])

			// dont make empty text nodes
			if len(copyText) != 0 {
				// we are textnode and has a "\n\n"
				p.AppendChild(NewText(copyText))
			}
			// dont add empty <p></p>
			if p.FirstChild != nil {
				block.Parent.InsertBefore(p, block)
				p = NewElement(tagName)
			}
		}

		if p.FirstChild != nil {
			block.Parent.InsertBefore(p, block)
		}
		block.Parent.RemoveChild(block)
	}
	return nil
}
