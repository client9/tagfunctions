package tagfunctions

import (
	"strings"

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

// GetKeyValue for shell style args
//
// Given a list of arguments in a format of "key=value"
// find the value that matches "key"
//
// This assume all keys do not or cannot have "=" in them.
func GetKeyValue(args []string, key string) (out string) {
	prefix := key + "="
	for _, kv := range args[1:] {
		if strings.HasPrefix(kv, prefix) {
			return kv[len(prefix):]
		}
	}
	return
}

// ToArgv converts a Node's attributes into shell "argv" style.
// Arg[0] is the name of the Node
// Arg[1...] are the attributes converted to the string "key=value"
//
//	if a value is empty then it's just "key"
func ToArgv(n *html.Node) []string {
	argv := make([]string, len(n.Attr)+1)
	argv[0] = n.Data
	for i, attr := range n.Attr {
		j := i + 1
		argv[j] = attr.Key
		if attr.Val != "" {
			argv[j] += "=" + attr.Val
		}
	}
	return argv
}

// SetArg (by index)
func SetArg(n *html.Node, i int, k string) {
	i--
	n.Attr[i].Key = k
	n.Attr[i].Val = ""
}

// GetArg (by index)
//
// 0 is Arg0, the name of the tag.
// 1 is the first attribute
//
// If index is invalid, an empty string is returned.
//
// TODO: can attr[i].value exit in this system?
//
//	if so, why not reutrn key=value
func GetArg(n *html.Node, i int) (out string) {
	if i == 0 {
		return n.Data
	}
	i--

	if i >= len(n.Attr) {
		return
	}
	return n.Attr[i].Key
}

// linear search cause it's likely so short
func getKey(n *html.Node, key string) (out string) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return
}
