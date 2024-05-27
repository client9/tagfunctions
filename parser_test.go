package tagfunctions

import (
	"strings"
	"testing"
)

func TestMacro2(t *testing.T) {

	type test struct {
		input string
		want  string
	}

	tests := []test{
		{"", "$root{}"},
		{"FOO", "$root{FOO}"},
		{"$", "$root{$}"},
		{"a$", "$root{a$}"},
		{"a$ ", "$root{a$ }"},
		{"$foo", "$root{$foo{}}"},
		{"abc$foo", "$root{abc$foo{}}"},
		{"$foo$bar", "$root{$foo{}$bar{}}"},
		{"$foo$bar next", "$root{$foo{}$bar{} next}"},
		{"$1.00", "$root{$1.00}"},
		{"$-1.00", "$root{$-1.00}"},
		{"$+1.00", "$root{$+1.00}"},
		{"$.00", "$root{$.00}"},
		{"$h1{headline}", "$root{$h1{headline}}"},
		{"1$h1{headline}2", "$root{1$h1{headline}2}"},
		{"$b{bold $i{italic}}", "$root{$b{bold $i{italic}}}"},
		{"$b{bold $i{italic} and bold again}", "$root{$b{bold $i{italic} and bold again}}"},
		{"plain $b{bold $i{italic}} text", "$root{plain $b{bold $i{italic}} text}"},
		{"$letters[a b third]", `$root{$letters[a b third]}`},
		{"$letters[  a b third   ]", `$root{$letters[a b third]}`},
		{"$letters[  a b third   ]after", `$root{$letters[a b third]after}`},
		{"$b[class ]{bold text}", `$root{$b[class]{bold text}}`},
		{"$b[class]{bold text}", `$root{$b[class]{bold text}}`},

		{"$b[c1 c2]{text}", `$root{$b[c1 c2]{text}}`},
		{"$b[  c1   c2  ]{text}", `$root{$b[c1 c2]{text}}`},
		{`$b[class=mega ]{boldx}`, `$root{$b[class=mega]{boldx}}`},
		{`$b[class=mega]{bold}`, `$root{$b[class=mega]{bold}}`},
		{`$b[class="mega"]{bold}`, `$root{$b[class=mega]{bold}}`},
		{`$b[class='mega']{bold}`, `$root{$b[class=mega]{bold}}`},
		{`$b[class='mega bold']{bold}`, `$root{$b[class="mega bold"]{bold}}`},
		{`$b[class="mega bold"]{bold}`, `$root{$b[class="mega bold"]{bold}}`},
		{`$b["1" "2" "3"]`, `$root{$b[1 2 3]}`},
		{`$b['1' '2' '3']`, `$root{$b[1 2 3]}`},

		// what happens when arg is quoted?  Rendering isn't 'correct' since attr key value
		// is actually invalid.
		//
		// see test case before for more explicit testing
		{`$b['class name']{bold}`, `$root{$b["class name"]{bold}}`},
		{`$b["class name"]{bold}`, `$root{$b["class name"]{bold}}`},
	}
	for i, tc := range tests {
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))
		sb := &strings.Builder{}

		if err := Render(sb, node); err != nil {
			t.Errorf("case %d: got unexpected error %v", i, err)
		}
		got := sb.String()
		if got != tc.want {
			t.Errorf("case %d: expected: %v, got %v", i, tc.want, got)
		}
	}
}

// what happen if an attribute name is quoted or has spaces?
/*
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
*/
