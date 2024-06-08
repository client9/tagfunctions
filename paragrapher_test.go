package tagfunctions

import (
	"strings"
	"testing"
)

// tests that existing paragraphs are split
func TestParagraphSplit(t *testing.T) {
	type test struct {
		input string
		want  string
	}
	tests := []test{
		{"", "<root></root>"},
		{"$p{line 1\nline 2}", "<root><p>line 1\nline 2</p></root>"},
		{"$p{line\n\n}", "<root><p>line</p></root>"},
		{"$p{line 1\n\nline 2}", "<root><p>line 1</p><p>line 2</p></root>"},
		{"$p{line 1   \n\n\n  line 2}", "<root><p>line 1</p><p>line 2</p></root>"},
		{"$p{$b{line 1}\n\n$b{line 2}}", "<root><p><b>line 1</b></p><p><b>line 2</b></p></root>"},
		{"$p{$b{bold}line1\n\n$b{line 2}}", "<root><p><b>bold</b>line1</p><p><b>line 2</b></p></root>"},
		{"$p{outer $p{inner} ending}", "<root><p>outer </p><p>inner</p><p> ending</p></root>"},
		{"$p{$pre{junk}outer}", "<root><pre>junk</pre><p>outer</p></root>"},
		{"$p{$pre{junk}\n\n}", "<root><pre>junk</pre></root>"},
		{"$p{line1$b{bold}\n$pre{junk}\nlast}", "<root><p>line1<b>bold</b>\n</p><pre>junk</pre><p>\nlast</p></root>"},
		{"$p{before $inline{middle} after}", "<root><p>before <inline>middle</inline> after</p></root>"},
	}
	pg := Paragrapher{}
	for num, tc := range tests {
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))
		if err := pg.Execute(node); err != nil {
			t.Fatalf("got Paragrapher error: %s", err)
		}
		sb := &strings.Builder{}
		if err := RenderHTML(sb, node); err != nil {
			t.Fatalf("case %d: got unexpected error %v", num, err)
		}
		got := sb.String()
		if tc.want != got {
			t.Errorf("Case %d: got %s want %s", num, got, tc.want)
		}
	}
}

func TestParagrapheRoot(t *testing.T) {
	type test struct {
		input string
		want  string
	}
	tests := []test{
		{"", "<root></root>"},
		{"line1\n\nline2", "<root><p>line1</p><p>line2</p></root>"},
		{"line1$p{line2}", "<root><p>line1</p><p>line2</p></root>"},
		{"line1$p{line2}line3", "<root><p>line1</p><p>line2</p><p>line3</p></root>"},
		{"$p{line1}line2", "<root><p>line1</p><p>line2</p></root>"},
		{"$p{line1}line2\n\nline3", "<root><p>line1</p><p>line2</p><p>line3</p></root>"},
		{"$p{line1}line2\n\n$p{line3}", "<root><p>line1</p><p>line2</p><p>line3</p></root>"},
	}

	pg := Paragrapher{}
	pg.Tag = "root"

	for num, tc := range tests {
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))
		if err := pg.Execute(node); err != nil {
			t.Fatalf("got Paragrapher error: %s", err)
		}
		sb := &strings.Builder{}
		if err := RenderHTML(sb, node); err != nil {
			t.Fatalf("case %d: got unexpected error %v", num, err)
		}
		got := sb.String()
		if tc.want != got {
			t.Errorf("Case %d: got %s want %s", num, got, tc.want)
		}
	}
}
