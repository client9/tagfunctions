package tagfunctions

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestTransformElement(t *testing.T) {
	n := &html.Node{
		Type: html.ElementNode,
		Data: "p",
	}

	sb := &strings.Builder{}
	TransformElement(n, "div", "class", "foo")
	RenderHTML(sb, n)

	t.Logf(sb.String())
}

func TestExecute(t *testing.T) {
	fmap := map[string]NodeFunc{
		"link": func(n *html.Node) error {
			TransformElement(n, "a")
			return nil
		},
	}
	n := NewElement("link", "href", "https://www.google.com/")
	sb := &strings.Builder{}
	if err := Execute(n, fmap); err != nil {
		t.Fatalf("failed: %s", err)
	}
	RenderHTML(sb, n)
	t.Logf(sb.String())
}

func TestExecuteFunc(t *testing.T) {

	fmap := map[string]NodeFunc{
		"link": func(n *html.Node) error {
			TransformElement(n, "a")
			return nil
		},
	}
	exec := ExecuteFunc(fmap)

	n := NewElement("link", "href", "https://www.google.com/")
	sb := &strings.Builder{}
	if err := exec(n); err != nil {
		t.Fatalf("failed: %s", err)
	}
	RenderHTML(sb, n)
	t.Logf(sb.String())

}
