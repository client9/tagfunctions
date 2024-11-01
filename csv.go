package tagfunctions

import (
	"encoding/csv"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// encodeForCSV takes a string and encodes it for safe embedding in a CSV file.
func encodeForCSV(s string) string {
	// Check if the string contains special characters
	if strings.ContainsAny(s, ",\"\n") {
		// Escape any double quotes by doubling them
		s = strings.ReplaceAll(s, `"`, `""`)
		// Surround the string with double quotes
		return `"` + s + `"`
	}
	// Return the string as-is if no special characters
	return s
}

func CsvEscape(n *html.Node) error {

	raw := ""
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.TextNode {
			continue
		}
		raw += child.Data
	}
	RemoveChildren(n)
	n.Type = html.TextNode
	n.Data = encodeForCSV(raw)
	n.DataAtom = 0
	n.Attr = nil
	return nil

}

// NewCsvTableHTML takes an embedded CSV and converts to an HTML table.
//
// The table's tags have optional class attributes using the formatter function.
// If nil, then no class attributes are added.
// It takes the class name, and the row and col of the cell if appropriate.
//
// fn("td", 3, 2") means emit a CSS class for a <td> in row 3, col 2.
// fn("table", 0,0,) the row and colum are always zero
//
// Sometimes you need to wrap a table in an outer div to get the desired behavior.
//
// fn("wrap", 0, 0)
//
// if non-empty will wrap the table in a <div class="xxx">
func NewCsvTableHTML(formatter func(string, int, int) string) NodeFunc {
	if formatter == nil {
		formatter = func(string, int, int) string { return "" }
	}
	makeTableTag := func(name string, row int, col int) *html.Node {
		cz := formatter(name, row, col)
		if cz != "" {
			return NewElement(name, "class", cz)
		}
		return NewElement(name)
	}

	return func(n *html.Node) error {
		body := n.FirstChild.Data

		table := makeTableTag("table", 0, 0)
		r := csv.NewReader(strings.NewReader(body))

		/* optional caption that probably needs work */
		/*
			caption := strings.Join(args[1:], " ")
			if caption != "" {
				captionElement := makeTableTag("caption", 0, 0)
				captionElement.AppendChild(NewText(caption))
				table.AppendChild(captionElement)
			}
		*/
		// read header row
		i := 0
		row, _ := r.Read()
		thead := makeTableTag("thead", 0, 0)
		tr := makeTableTag("tr", i, 0)
		table.AppendChild(thead)
		thead.AppendChild(tr)
		for j, col := range row {
			th := makeTableTag("th", i, j)
			th.AppendChild(NewText(col))
			thead.AppendChild(th)
		}
		tbody := makeTableTag("tbody", 0, 0)
		table.AppendChild(tbody)

		for {
			i += 1
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
				log.Printf("row %v: table logger: %v", row, err)
				break
			}
			tr = makeTableTag("tr", i, 0)
			tbody.AppendChild(tr)
			for j, col := range row {
				td := makeTableTag("td", i, j)
				td.AppendChild(NewText(col))
				tr.AppendChild(td)
			}
		}

		// table is complete.  Link it in.
		// maybe wrap it
		wrap := formatter("wrapper", 0, 0)
		if wrap != "" {
			tmp := NewElement("div", "class", "wrap")
			tmp.AppendChild(table)
			table = tmp
		}

		// Replace node
		n.Parent.InsertBefore(table, n)
		n.Parent.RemoveChild(n)
		return nil
	}
}
