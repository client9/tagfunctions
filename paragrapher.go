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

type Paragrapher struct {
	Tag      string // elements to split on "\n\n" to generate new blocks
	Create   string // element to make
	IsInline func(*html.Node) bool
}

func (p *Paragrapher) Execute(n *html.Node) error {

	//
	// Set defaults
	//
	if p.Tag == "" {
		p.Tag = "p"
	}
	if p.Create == "" {
		if p.Tag == "root" {
			p.Create = "p"
		} else {
			p.Create = p.Tag
		}
	}

	if p.IsInline == nil {
		p.IsInline = inlineNode
	}

	if err := p.executeTag(n, p.Tag); err != nil {
		return err
	}

	return nil
}

func inlineNode(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return true
	}
	switch n.Data {

	// are we standard HTML inline tags?
	case "a", "b", "em", "i", "span", "sup", "sub", "small", "code", "addr", "strong", "var", "label", "cite", "tt", "kbd", "time", "br", "big", "acronym":
		return true

	// are we standard HTML block tags?
	case "hr", "root", "div", "p", "pre", "blockquote", "article", "section", "table", "img", "figure", "h1", "h2", "h3", "h4", "h5", "h6":

		return false
	}

	// Are we a tag with no children?
	if n.FirstChild == nil {
		return true
	}

	// Are we a tag with exactly one child that is text node?
	// Then assume it's inline
	if n.FirstChild != nil && n.FirstChild.NextSibling == nil && n.FirstChild.Type == html.TextNode {
		return true
	}

	return false
}

func needsSplit(n *html.Node) bool {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.TextNode {
			continue
		}
		if strings.Contains(child.Data, "\n\n") {
			return true
		}
	}
	return false
}

func (pg *Paragrapher) executeTag(n *html.Node, tagName string) error {
	blocks := Selector(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagName
	})

	//
	tagName = pg.Create

	for _, block := range blocks {

		// possible we don't need to do all this cut and paste

		p := NewElement(tagName)
		current := block.FirstChild
		for current != nil {
			if !pg.IsInline(current) {
				if p.FirstChild != nil {
					if block.Parent != nil {
						block.Parent.InsertBefore(p, block)
					} else {
						block.InsertBefore(p, current)
					}
					p = NewElement(tagName)
				}
				next := current.NextSibling
				if block.Parent != nil {
					block.RemoveChild(current)
					block.Parent.InsertBefore(current, block)
				}
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

			copyText := strings.TrimRight(current.Data[:idx], " \n\r\t")
			current.Data = strings.TrimLeft(current.Data[idx+2:], " \n\r\t")

			// dont make empty text nodes
			if len(copyText) != 0 {
				// we are textnode and has a "\n\n"
				p.AppendChild(NewText(copyText))
			}
			// dont add empty <p></p>
			if p.FirstChild != nil {
				if block.Parent != nil {
					block.Parent.InsertBefore(p, block)
				} else {
					block.InsertBefore(p, current)
				}
				p = NewElement(tagName)
			}
		}

		if p.FirstChild != nil {
			if block.Parent != nil {
				block.Parent.InsertBefore(p, block)
			} else {
				block.InsertBefore(p, current)
			}
		}
		if block.Parent != nil {
			block.Parent.RemoveChild(block)
		}
	}
	return nil
}
