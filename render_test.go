package tagfunctions

import (
	"strings"
	"testing"
)

// tests identity  orig -> parse -> render -> orig
// and HTML variation
func TestRender(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{"", ""},
		{"$b{bold} text", "<b>bold</b> text"},
		{"$b{bold $i{italic} text}", "<b>bold <i>italic</i> text</b>"},
		{"$p[class=text]{body}", `<p class="text">body</p>`},
		{"$echo[1 2 3 4]", `<echo 1="" 2="" 3="" 4=""></echo>`},
	}
	for i, tc := range tests {

		// parse
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))

		// test that is renders back to original
		sb := strings.Builder{}
		Render(&sb, node)
		got := sb.String()
		want := "$root{" + tc.input + "}"
		if got != want {
			t.Errorf("case %d: expected: %v, got %v", i, want, got)
		}

		// test HTML rendering
		sb.Reset()
		RenderHTML(&sb, node)
		got = sb.String()
		want = "<root>" + tc.want + "</root>"
		if got != want {
			t.Errorf("case %d: expected: %v, got %v", i, want, got)
		}
	}
}
