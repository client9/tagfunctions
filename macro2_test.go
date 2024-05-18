package tagfunctions

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestMacro2(t *testing.T) {

	type test struct {
		input string
		want  string
	}

	tests := []test{
		{"", "<root></root>"},
		{"FOO", "<root>FOO</root>"},
		{"$", "<root>$</root>"},
		{"a$", "<root>a$</root>"},
		{"a$ ", "<root>a$ </root>"},
		{"$foo", "<root><foo></foo></root>"},
		{"abc$foo", "<root>abc<foo></foo></root>"},
		{"$foo$bar", "<root><foo></foo><bar></bar></root>"},
		{"$foo$bar next", "<root><foo></foo><bar></bar> next</root>"},
		{"$1.00", "<root>$1.00</root>"},
		{"$-1.00", "<root>$-1.00</root>"},
		{"$+1.00", "<root>$+1.00</root>"},
		{"$.00", "<root>$.00</root>"},
		{"$h1{headline}", "<root><h1>headline</h1></root>"},
		{"1$h1{headline}2", "<root>1<h1>headline</h1>2</root>"},
		{"$b{bold $i{italic}}", "<root><b>bold <i>italic</i></b></root>"},
		{"$b{bold $i{italic} and bold again}", "<root><b>bold <i>italic</i> and bold again</b></root>"},
		{"plain $b{bold $i{italic}} text", "<root>plain <b>bold <i>italic</i></b> text</root>"},
		{"$letters[a b third]", `<root><letters a="" b="" third=""></letters></root>`},
		{"$letters[  a b third   ]", `<root><letters a="" b="" third=""></letters></root>`},
		{"$letters[  a b third   ]after", `<root><letters a="" b="" third=""></letters>after</root>`},
		{"$b[class ]{bold text}", `<root><b class="">bold text</b></root>`},
		{"$b[class]{bold text}", `<root><b class="">bold text</b></root>`},
		{"$b[c1 c2]{text}", `<root><b c1="" c2="">text</b></root>`},
		{"$b[  c1   c2  ]{text}", `<root><b c1="" c2="">text</b></root>`},
		{`$b[class=mega ]{boldx}`, `<root><b class="mega">boldx</b></root>`},
		{`$b[class=mega]{bold}`, `<root><b class="mega">bold</b></root>`},
		{`$b[class="mega"]{bold}`, `<root><b class="mega">bold</b></root>`},
		{`$b[class='mega']{bold}`, `<root><b class="mega">bold</b></root>`},
		{`$b[class='mega bold']{bold}`, `<root><b class="mega bold">bold</b></root>`},

		// what happens when arg is quoted?  Rendering isn't 'correct' since attr key value
		// is actually invalid.
		//
		// see test case before for more explicit testing
		{`$b['class name']{bold}`, `<root><b class name="">bold</b></root>`},
	}
	for i, tc := range tests {
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))
		sb := &strings.Builder{}

		if err := html.Render(sb, node); err != nil {
			t.Fatalf("case %d: got unexpected error %v", i, err)
		}
		got := sb.String()
		if got != tc.want {
			t.Fatalf("case %d: expected: %v, got %v", i, tc.want, got)
		}
	}
}

// what happen if an attribute name is quoted or has spaces?
func TestQuotedAttributeName(t *testing.T) {
	input := `$b['mega bold']{hello}`
	p := Tokenizer{}
	root := p.Parse(strings.NewReader(input))
	if root.Data != "root" {
		t.Fatalf("root node is %s", root.Data)
	}
	child := root.FirstChild
	if child == nil || child.Data != "b" {
		t.Fatalf("child node is %v", child)
	}
	attr := child.Attr
	if len(attr) != 1 {
		t.Fatalf("expected 1 child attribute, got %d", len(attr))
	}
	if attr[0].Key != "class name" && attr[0].Val != "" {
		t.Fatalf("expect key:'class name' got: %v", attr[0])
	}
}
