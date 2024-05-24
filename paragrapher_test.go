package tagfunctions

import (
	"strings"
	"testing"
)

func TestParagraph(t *testing.T) {
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
