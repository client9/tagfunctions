package tagfunctions

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type test struct {
		input string
		args  []string
		body  string
	}

	tests := []test{
		{
			input: "name John Brown",
			args:  []string{"name", "John", "Brown"},
		},
		{
			input: "name 'John' \"Brown\"",
			args:  []string{"name", "John", "Brown"},
		},
		{
			input: "name John Brown\n\n\n",
			args:  []string{"name", "John", "Brown"},
		},
		{
			input: "   name   John   Brown   \n\n\n",
			args:  []string{"name", "John", "Brown"},
		},
		{
			// check to make sure flags are ok
			input: "name -first John -last Brown",
			args:  []string{"name", "-first", "John", "-last", "Brown"},
		},
		{
			// check to make sure key=val are ok
			input: "name first=John last=Brown",
			args:  []string{"name", "first=John", "last=Brown"},
		},
		{
			input: "name John Brown {\napple {banana}}",
			args:  []string{"name", "John", "Brown"},
			body:  "\napple {banana}",
		},
	}

	for i, tc := range tests {
		s := NewScanner([]byte(tc.input))
		args, body, err := s.Next()
		if err != nil {
			t.Fatalf("case %d, got error %v", i, err)
		}
		if !reflect.DeepEqual(args, tc.args) {
			t.Fatalf("case %d, expected args %v got %v", i, tc.args, args)
		}
		if tc.body != body {
			t.Fatalf("case %d, expected body %q got %q", i, tc.body, body)
		}
	}
}
