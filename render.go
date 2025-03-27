package tagfunctions

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func RenderStringFunc(render func(io.Writer, *html.Node) error) func(n *html.Node) ([]byte, error) {
	return func(n *html.Node) ([]byte, error) {
		buf := &bytes.Buffer{}
		if err := render(buf, n); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

// Render AST into HTML
func RenderHTML(w io.Writer, n *html.Node) error {
	// not much to do here!
	return html.Render(w, n)
}

// hack used in many places in go
type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

// Render AST back into parseble string
//
//	i.e. $root{....}
func Render(w io.Writer, n *html.Node) error {

	// hack used many places in golang
	// if the Writer is say... StringWriter
	// then nothing to do..
	// othewise wrap it in a bufio.Writer
	//
	if x, ok := w.(writer); ok {
		return render1(x, n)
	}
	buf := bufio.NewWriter(w)
	if err := render1(buf, n); err != nil {
		return err
	}
	return buf.Flush()
}

// render attribute
// while are reusing the html.Node and html.Attribute
// the attributes here have no HTML restrictions.
func renderAttr(attr html.Attribute) string {
	k := attr.Key
	v := attr.Val

	simpleKey := !strings.ContainsAny(k, "\n\r '\"")
	simpleValue := !strings.ContainsAny(v, "\n\r '\"")

	// key
	// key=value
	// key="value"
	if simpleKey {
		if v == "" {
			return k
		}
		if simpleValue {
			return k + "=" + v
		}
		return fmt.Sprintf("%s=%q", k, v)
	}
	// key is value (ok)
	if v == "" {
		return fmt.Sprintf("%q", k)
	}
	// this case should never happen but it's ok.
	return fmt.Sprintf("%q", k+"="+v)
}

func render1(w writer, n *html.Node) error {
	// Render non-element nodes; these are the easy cases.
	switch n.Type {
	case html.ErrorNode:
		return errors.New("render: an ErrorNode node")
	case html.TextNode:
		_, err := w.WriteString(n.Data)
		return err
	case html.ElementNode:
		//NOP
	default:
		return errors.New("render: unknown node type")
	}

	w.WriteByte('$')
	w.WriteString(n.Data)
	if len(n.Attr) > 0 {
		args := make([]string, len(n.Attr))
		for i, a := range n.Attr {
			args[i] = renderAttr(a)
		}
		w.WriteString("[" + strings.Join(args, " ") + "]")

		// no children?  We are done
		if n.FirstChild == nil {
			return nil
		}
	}

	// render children
	w.WriteByte('{')
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := render1(w, c); err != nil {
			return err
		}
	}
	return w.WriteByte('}')
}
