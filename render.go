package tagfunctions

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// takes a HTML-node tree and renders it using user functions.
// some node might be "pass through" (i.e. just render back to HTML).
//
// could be cleaned up but this is mostly simple and mostly happy with it.
//

type TagFunc func([]string, string) string

// MakeTag creates a simple HTML tag with no attributes.
//
// TODO: why not pass through arguments?
func MakeTag(tag string) TagFunc {
	return func(args []string, body string) string {
		return fmt.Sprintf("<%s>%s</%s>", tag, strings.TrimSpace(body), tag)
	}
}

// MakeTagClass makes a HTML with a class attribute
func MakeTagClass(tag string, cz string) TagFunc {
	return func(args []string, body string) string {
		return fmt.Sprintf("<%s class=%q>%s</%s>", tag, cz, strings.TrimSpace(body), tag)
	}
}

func IsHTMLTag(tag string) bool {
	// todo make map
	switch tag {
	case "sup":
		return true
	case "h1", "h2", "h3", "h4", "p", "b", "i", "em", "hr", "s":
		return true
	case "table", "tbody", "th", "tr", "td", "tfoot":
		return true
	case "blockquote", "pre":
		return true
	case "ul", "ol", "li":
		return true
	}
	return false
}
func HTMLTag(n *html.Node, body string) string {
	// hack for now
	out := "<" + n.Data
	for _, a := range n.Attr {
		out += " " + a.Key + "=" + fmt.Sprintf("%q", a.Val)
	}
	out += fmt.Sprintf(">%s</%s>", body, n.Data)
	return out
}

func RenderFunc(fmap map[string]TagFunc) func(n *html.Node) string {
	return func(n *html.Node) string {
		return Render(n, fmap)
	}
}

// "select title" --
//
//	  only returns text children
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

func Render(n *html.Node, fmap map[string]TagFunc) string {
	switch n.Type {
	case html.TextNode:
		return n.Data
	case html.ElementNode:
		body := ""
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			body += Render(c, fmap)
		}
		if fn, ok := fmap[n.Data]; ok {
			out := fn(toArgv(n), body)
			return out
		}
		if IsHTMLTag(n.Data) {
			return HTMLTag(n, body)
		}
		// unknown tag ... pass through
		out := "$" + n.Data
		if len(n.Attr) > 0 {
			out += "["
			htmlargs := []string{}
			for _, arg := range n.Attr {
				if len(arg.Val) == 0 {
					// TODO: better quoting
					if strings.Contains(arg.Key, " ") {
						htmlargs = append(htmlargs, fmt.Sprintf("%q", arg.Key))
						continue
					}
					htmlargs = append(htmlargs, arg.Key)
					continue
				}
				// TODO: better quoting
				if strings.Contains(arg.Val, " ") {
					htmlargs = append(htmlargs, fmt.Sprintf("%s=%q", arg.Key, arg.Val))
					continue
				}
				htmlargs = append(htmlargs, arg.Key+"="+arg.Val)
			}
			out += strings.Join(htmlargs, " ")
			out += "]"
		}
		body = strings.TrimSpace(body)
		if len(body) > 0 {
			out += "{" + body + "}"
		}
		return "<code>" + out + "</code>"
	default:
		panic("unknown node type")
	}
}
