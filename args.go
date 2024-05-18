package tagfunctions

import (
	"strings"

	"golang.org/x/net/html"
)

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

// for shell style args
func getKeyValue(args []string, key string) (out string) {
	prefix := key + "="
	for _, kv := range args[1:] {
		if strings.HasPrefix(kv, prefix) {
			return kv[len(prefix):]
		}
	}
	return
}

// Unix shell style argv.
// argv[0] is name of the node, followed by arguments
func toArgv(n *html.Node) []string {
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

func setArg(n *html.Node, i int, k string) {
	i--
	n.Attr[i].Key = k
	n.Attr[i].Val = ""
}

// no error checking cause that can be done a different way
func getArg(n *html.Node, i int) (out string) {
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
