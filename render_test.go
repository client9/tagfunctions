package tagfunctions

import (
	"strings"
	"testing"
)

func TestRedner(t *testing.T) {
	fmap := map[string]TagFunc{
		"root": func(argv []string, body string) string {
			return body
		},
		"b": func(argv []string, body string) string {
			return "<b>" + body + "</b>"
		},
		"i": func(argv []string, body string) string {
			return "<i>" + body + "</i>"
		},
		"p": func(argv []string, body string) string {
			cname := GetKeyValue(argv, "name")
			return "<p class=" + cname + ">" + body + "</p>"
		},
		"echo": func(argv []string, body string) string {
			return "<echo>" + strings.Join(argv[1:], " ") + "</echo>"
		},
	}

	type test struct {
		input string
		want  string
	}

	tests := []test{
		{"", ""},
		{"$b{bold} text", "<b>bold</b> text"},
		{"$b{bold $i{italic} text}", "<b>bold <i>italic</i> text</b>"},
		{"$p[name=text]{body}", "<p class=text>body</p>"},
		{"$echo[1 2 3 4]", "<echo>1 2 3 4</echo>"},
	}
	for i, tc := range tests {
		p := Tokenizer{}
		node := p.Parse(strings.NewReader(tc.input))
		got := Render(node, fmap)
		if got != tc.want {
			t.Fatalf("case %d: expected: %v, got %v", i, tc.want, got)
		}
	}
}
