package tagfunctions

import (
	"io"
	"log"

	"golang.org/x/net/html"
)

type Tokenizer struct {
	r         io.ByteScanner
	maybeText []byte
	current   *html.Node
}

func (z *Tokenizer) readByte() (byte, error) {
	return z.r.ReadByte()
}
func (z *Tokenizer) unreadByte() {
	if err := z.r.UnreadByte(); err != nil {
		// should never happen
		panic("asset failed: unread byte failed")
	}
}

func (z *Tokenizer) Parse(r io.ByteScanner) *html.Node {
	z.maybeText = nil
	z.r = r
	root := &html.Node{
		Type: html.ElementNode,
		Data: "root",
	}
	z.current = root
	z.stateText()
	return root
}

func (z *Tokenizer) stateText() {
	for {
		c, err := z.readByte()
		if err != nil {
			if len(z.maybeText) == 0 {
				return
			}
			// append final text node
			text := &html.Node{
				Type: html.TextNode,
				Data: string(z.maybeText),
			}
			if z.current == nil {
				panic("current is nil: adding " + text.Data)
			}
			z.current.AppendChild(text)
			z.maybeText = nil
			return
		}
		switch c {
		case '$':
			z.stateAfterDollar()
		case '}':
			if len(z.maybeText) > 0 {
				// append final text node
				text := &html.Node{
					Type: html.TextNode,
					Data: string(z.maybeText),
				}
				z.current.AppendChild(text)
				z.maybeText = nil
			}
			if z.current.Parent != nil {
				z.current = z.current.Parent
			}
		default:
			z.maybeText = append(z.maybeText, c)
		}
	}
}

func (z *Tokenizer) stateAfterDollar() {
	for {
		c, err := z.readByte()
		if err != nil {
			z.maybeText = append(z.maybeText, '$')
			// append final text node
			text := &html.Node{
				Type: html.TextNode,
				Data: string(z.maybeText),
			}
			z.current.AppendChild(text)
			z.maybeText = nil
			return
		}
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.', ',', ' ', '\t', '\r', '\f', '\n':
			z.maybeText = append(z.maybeText, '$')
			z.maybeText = append(z.maybeText, c)
			return
		case '[':
			// TBD
		case '{':
			log.Fatalf("Got '${', expected a function name")
		default:
			if len(z.maybeText) > 0 {
				text := &html.Node{
					Type: html.TextNode,
					Data: string(z.maybeText),
				}
				z.current.AppendChild(text)
				z.maybeText = nil
			}
			z.stateFunctionName(c)
			return
		}
	}
}

func (z *Tokenizer) stateFunctionName(c byte) {
	fname := []byte{c}
	for {
		c, err := z.readByte()
		if err != nil {
			// $x is valid.. attach node
			n := &html.Node{
				Type: html.ElementNode,
				Data: string(fname),
			}
			z.current.AppendChild(n)
			return
		}

		switch c {
		// case ']':
		// probably an error
		case '[':
			// $foo[.... start of args.  Assume valid node
			n := &html.Node{
				Type: html.ElementNode,
				Data: string(fname),
			}
			z.stateBeforeAttributeName(n)
			return
		case ' ', '\t', '\r', '\f', '\n':
			// $FOO
			n := &html.Node{
				Type: html.ElementNode,
				Data: string(fname),
			}
			z.current.AppendChild(n)
			z.unreadByte()
			return
		case '$':
			// $FOO$BAR
			n := &html.Node{
				Type: html.ElementNode,
				Data: string(fname),
			}
			z.current.AppendChild(n)
			z.unreadByte()
			return
		case '{':
			n := &html.Node{
				Type: html.ElementNode,
				Data: string(fname),
			}
			z.current.AppendChild(n)
			z.current = n
			z.stateText()
			return
		default:
			fname = append(fname, c)
		}
	}
}

func (z *Tokenizer) stateBeforeAttributeName(n *html.Node) {
	for {
		c, err := z.readByte()
		if err != nil {
			// TBD on what to do here
			return
		}

		switch c {
		case ' ', '\t', '\f', '\r', '\n':
			continue
		case ']':
			// $foo[]....
			//
			z.stateAfterAttributes(n)
			return
		case '\'':
			z.stateAttributeNameQuote1(n)
		case '"':
			z.stateAttributeNameQuote2(n)
		default:
			z.unreadByte()
			z.stateAttributeName(n)
		}
	}
}

func (z *Tokenizer) stateAttributeNameQuote1(n *html.Node) {
	key := []byte{}
	for {
		c, err := z.readByte()
		if err != nil {
			// TBD on what to do here
			return
		}

		switch c {
		case '\'':
			n.Attr = append(n.Attr, html.Attribute{Key: string(key)})
			z.stateAfterAttributeValueQuoted()
			return
		default:
			key = append(key, c)
		}
	}
}
func (z *Tokenizer) stateAttributeNameQuote2(n *html.Node) {
	key := []byte{}
	for {
		c, err := z.readByte()
		if err != nil {
			// TBD on what to do here
			return
		}

		switch c {
		case '"':
			n.Attr = append(n.Attr, html.Attribute{Key: string(key)})
			z.stateAfterAttributeValueQuoted()
			return
		default:
			key = append(key, c)
		}
	}
}

func (z *Tokenizer) stateAttributeName(n *html.Node) {
	key := []byte{}
	for {
		c, err := z.readByte()
		if err != nil {
			// TBD on what to do here
			return
		}

		switch c {
		case ' ', '\t', '\f', '\r', '\n':
			// $foo[xxxi<sp>
			n.Attr = append(n.Attr, html.Attribute{Key: string(key)})
			return
		case ']':
			// $foo[]....
			//
			n.Attr = append(n.Attr, html.Attribute{Key: string(key)})
			z.unreadByte()
			return
		case '=':
			z.stateBeforeAttributeValue(n, string(key))
			return
		default:
			key = append(key, c)
		}
	}
}

// stateBeforeAttributeValue
// https://html.spec.whatwg.org/#before-attribute-value-state
//
//	key=value or key= value key="value"  key='value'
func (z *Tokenizer) stateBeforeAttributeValue(n *html.Node, key string) {
	for {
		c, err := z.readByte()
		if err != nil {
			// ERROR
			return
		}

		switch c {
		case ' ', '\t', '\f', '\r', '\n':
			continue
		case '\'':
			z.stateAttributeValueQuote1(n, key)
			return
		case '"':
			z.stateAttributeValueQuote2(n, key)
			return
		case ']':
			// key=]
			// TBD
		default:
			z.unreadByte()
			z.stateAttributeValue(n, key)
			return
		}
	}
}

// state Attribute Value, Single Quoted
// https://html.spec.whatwg.org/#attribute-value-(single-quoted)-state
func (z *Tokenizer) stateAttributeValueQuote1(n *html.Node, key string) {
	val := []byte{}

	for {
		c, err := z.readByte()
		if err != nil {
			// ERROR
			return
		}

		switch c {
		case '\'':
			a := html.Attribute{Key: key, Val: string(val)}
			n.Attr = append(n.Attr, a)
			z.stateAfterAttributeValueQuoted()
			return
		default:
			val = append(val, c)
		}
	}
}

// state Attribute Value, Double Quoted
// https://html.spec.whatwg.org/#attribute-value-(double-quoted)-state
func (z *Tokenizer) stateAttributeValueQuote2(n *html.Node, key string) {
	val := []byte{}

	for {
		c, err := z.readByte()
		if err != nil {
			// ERROR
			return
		}

		switch c {
		case '"':
			a := html.Attribute{Key: key, Val: string(val)}
			n.Attr = append(n.Attr, a)
			z.stateAfterAttributeValueQuoted()
			return
		default:
			val = append(val, c)
		}
	}
}

func (z *Tokenizer) stateAfterAttributeValueQuoted() {
	c, err := z.readByte()
	if err != nil {
		// ERRROR
		return
	}

	switch c {
	case ' ', '\t', '\f', '\n', '\r':
		return
	case ']':
		z.unreadByte()
		return
	default:
		// ERROR.. only whitespace or ']' after quote
		// e.g. foo='bar'xxx
		return
	}
}

func (z *Tokenizer) stateAttributeValue(n *html.Node, key string) {
	val := []byte{}
	for {
		c, err := z.readByte()
		if err != nil {
			// ERROR
			return
		}

		switch c {
		case ' ', '\t', '\f', '\r', '\n':
			n.Attr = append(n.Attr, html.Attribute{Key: key, Val: string(val)})
			return
		case ']':
			n.Attr = append(n.Attr, html.Attribute{Key: key, Val: string(val)})
			z.unreadByte()
			return
		case '\'':
			// note quite right
			// this is the case of
			//  key=foo'bar'
			z.stateAttributeValueQuote1(n, key)
			return
		case '"':
			// not quite right
			// this is the case of
			//   key=foo"bar"
			//
			z.stateAttributeValueQuote2(n, key)
			return
		default:
			val = append(val, c)
		}
	}
}
func (z *Tokenizer) stateAfterAttributes(n *html.Node) {
	z.current.AppendChild(n)

	c, err := z.readByte()
	if err != nil {
		// exactly $foo[...]<EOF>
		return
	}
	switch c {
	case '{':
		z.current = n
		return
	}
	z.unreadByte()
}
