package tagfunctions

import (
	"log"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Execute parses and renders a tag string
func Generate(src string, fmap map[string]NodeFunc) (string, error) {
	t := Tokenizer{}
	n := t.Parse(strings.NewReader(src))
	if err := Execute(n, fmap); err != nil {
		return "", err
	}
	s := &strings.Builder{}
	if err := Render(s, n); err != nil {
		return "", err
	}
	return s.String(), nil
}

// Append adds the child to the parent and returns the parent
// helper to do function chaining
func Append(parent *html.Node, children ...*html.Node) *html.Node {
	for _, child := range children {
		parent.AppendChild(child)
	}
	return parent
}

// Replace
func Replace(dst, src *html.Node) *html.Node {

	RemoveChildren(dst)

	// move children over
	Reparent(dst, src)

	// copy attributes
	dst.Attr = make([]html.Attribute, len(src.Attr))
	copy(dst.Attr, src.Attr)

	// copy basics
	dst.Type = src.Type
	dst.Data = src.Data
	dst.DataAtom = src.DataAtom

	return dst
}

// RemoveChildren removes all child nodes and returns itself
func RemoveChildren(n *html.Node) *html.Node {
	for n.FirstChild != nil {
		n.RemoveChild(n.FirstChild)
	}
	return n
}

// Reparent moves children from src to dst, and returns dst
func Reparent(dst, src *html.Node) *html.Node {
	for {
		child := src.FirstChild
		if child == nil {
			break
		}
		src.RemoveChild(child)
		dst.AppendChild(child)
	}
	return dst
}

func NewElement(name string, kv ...string) *html.Node {
	// TODO ATOMS

	n := &html.Node{
		Type: html.ElementNode,
		Data: name,
	}
	for i := 0; i < len(kv); i += 2 {
		n.Attr = append(n.Attr, html.Attribute{Key: kv[i], Val: kv[i+1]})
	}
	return n
}

func NewText(text string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: text,
	}
}

// takes a HTML-node tree and renders it using user functions.
// some node might be "pass through" (i.e. just render back to HTML).
//
// could be cleaned up but this is mostly simple and mostly happy with it.
//

type NodeFunc func(n *html.Node) error

// MakeTag return a NodeFunc that transforms the incoming node
// - Type of ElementNode
// - With a new Tag name
// - With no Attributes
// - Children are preserved
func MakeTag(tag string) NodeFunc {
	return func(n *html.Node) error {
		TransformElement(n, "p")
		return nil
	}
}

// MakeTagClass returns a NodeFunc that transforms in incoming node:
// - Type of ElementNode
// - With a new Tag Name
// - Clears all existing attributes
// - Add an attribute of class with value
// - Children are preserved.
func MakeTagClass(tag string, cz string) NodeFunc {
	return func(n *html.Node) error {
		TransformElement(n, tag, "class", cz)
		return nil
	}
}

/*
func RenderFunc(fmap map[string]NodeFunc) func(n *html.Node) string {
	return func(n *html.Node) string {
		return Execute(n, fmap)
	}
}
*/

// "select title" --
//
//	only returns text children
//	Could be improved.
func Select(n *html.Node, tag string) string {
	blocks := Selector(n, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tag
	})
	if len(blocks) == 0 {
		return "nope"
	}
	body := ""
	for c := blocks[0].FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			body += c.Data
		}
	}
	return body
}

// Transform changes ElementNode's name and attributes
// children remain the same
func TransformElement(n *html.Node, name string, attr ...string) *html.Node {
	if n.Type != html.ElementNode {
		panic("not an element node")
	}
	if name == "" {
		panic("changing an element node to no-name")
	}
	if len(attr)&1 == 1 {
		panic("odd number of args given")
	}
	n.DataAtom = atom.Lookup([]byte(name))
	n.Data = name
	n.Attr = nil
	if len(attr) == 0 {
		return n
	}
	n.Attr = make([]html.Attribute, len(attr)/2)
	j := 0
	for i := 0; i < len(attr); i += 2 {
		n.Attr[j] = html.Attribute{Key: attr[i], Val: attr[i+1]}
		j++
	}
	return n
}

func TextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	out := ""
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		switch n.Type {
		case html.TextNode:
			out += TextContent(n)
		case html.ElementNode:
			out += TextContent(n)
		default:
			// NOP
		}
	}
	return out
}

func ExecuteFunc(m map[string]NodeFunc) NodeFunc {
	return func(n *html.Node) error {
		return Execute(n, m)
	}
}

func Execute(n *html.Node, fmap map[string]NodeFunc) error {
	switch n.Type {
	case html.TextNode:
		// TODO: Text processing function
		return nil
	case html.ElementNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := Execute(c, fmap); err != nil {
				return err
			}
		}
		log.Printf("NODE: %s", n.Data)
		if fn, ok := fmap[n.Data]; ok {
			return fn(n)
		}
		//
	default:
		panic("unknown node type")
	}
	return nil
}
