package tagfunctions

import (
	"golang.org/x/net/html"
)

// Selector recursively selects all nodes that match a given function
//
// TODO: why is this in this file?
func Selector(n *html.Node, fn func(*html.Node) bool) []*html.Node {
	out := []*html.Node{}
	if fn(n) {
		out = append(out, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		out = append(out, Selector(c, fn)...)
	}
	return out
}

// ToArgs converts a Node's attributes into a list of arguments
//
// if a value is empty then it's just "key"
func ToArgs(n *html.Node) []string {
	if len(n.Attr) == 0 {
		return nil
	}
	argv := make([]string, len(n.Attr))
	for i, attr := range n.Attr {
		argv[i] = attr.Key
		if attr.Val != "" {
			argv[i] += "=" + attr.Val
		}
	}
	return argv
}

// SetArg Set attribute argument by index
func SetArg(n *html.Node, i int, k string) {
	n.Attr[i].Key = k
	n.Attr[i].Val = ""
}

// GetArg - get Attribute value by index.
//
// If index is invalid, an empty string is returned.
func GetArg(n *html.Node, i int) (out string) {
	if i >= len(n.Attr) {
		return
	}
	key := n.Attr[i].Key
	val := n.Attr[i].Val
	if val == "" {
		// this should be the most common case
		return key
	}
	return key + "=" + val
}

// GetAttr - Get Attribute by key
func GetAttr(n *html.Node, key string) (out string) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return
}
