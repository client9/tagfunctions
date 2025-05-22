package tagfunctions

import (
	"log"
	"testing"
)

func TestTable1(t *testing.T) {
	fmap := make(map[string]NodeFunc)
	fmap["csvtable"] = NewCsvTableHTML(nil)

	doc := `$csvtable{
A,B,C
1,2,3
}`
	got, err := Generate(doc, fmap)
	if err != nil {
		t.Fatalf("Failed: %s", err)
	}
	want := "$root{$table{$thead{$tr{$th{A}$th{B}$th{C}}}$tbody{$tr{$td{1}$td{2}$td{3}}}}}"
	if want != got {
		t.Fatalf("Failed want %q, got %q", want, got)
	}

	doc2 := `$csvtable{
A,B,Cdumb
1,2,$b[class='bluebig']{3}
}`

	got, err = Generate(doc, fmap)
	if err != nil {
		t.Fatalf("Generate Plain Failed: %s", err)
	}
	log.Printf("Generate HTML: %s", got)
	got, err = GenerateHTML(doc2, fmap)
	if err != nil {
		t.Fatalf("Generate HTML Failed: %s", err)
	}
	log.Printf("FINAL: %s", got)
}
