package tagfunctions

import (
	"encoding/csv"
	"strings"
)

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
func NewCsvTableHTML(formatter func(string, int, int) string) TagFunc {
	if formatter == nil {
		formatter = func(string, int, int) string { return "" }
	}
	makeTableTag := func(name string, row int, col int) string {
		cz := formatter(name, row, col)
		if cz != "" {
			// TODO: escape attribute value
			cz = " class='" + cz + "'"
		}
		return "<" + name + cz + ">"
	}

	return func(args []string, body string) string {
		out := strings.Builder{}

		wrap := formatter("wrapper", 0, 0)
		if wrap != "" {
			// TODO: escape wrap
			out.WriteString("<div class=" + wrap + ">")
		}
		out.WriteString(makeTableTag("table", 0, 0))
		r := csv.NewReader(strings.NewReader(body))

		/* optional caption that probably needs work */
		caption := strings.Join(args[1:], " ")
		if caption != "" {
			out.WriteString(makeTableTag("caption", 0, 0))
			out.WriteString(caption)
			out.WriteString("</caption>")
		}
		// read header row
		i := 0
		row, _ := r.Read()
		out.WriteString(makeTableTag("thead", 0, 0))
		out.WriteString(makeTableTag("tr", i, 0))
		for j, col := range row {
			out.WriteString(makeTableTag("th", i, j))
			out.WriteString(col)
			out.WriteString("</th>")
		}
		out.WriteString("</tr></thead>\n")
		out.WriteString(makeTableTag("tbody", 0, 0))
		for {
			i += 1
			row, err := r.Read()
			if err != nil {
				break
			}
			out.WriteString(makeTableTag("tr", i, 0))
			for j, col := range row {
				out.WriteString(makeTableTag("td", i, j))
				out.WriteString(col)
				out.WriteString("</td>")
			}
			out.WriteString("</tr>\n")
		}
		out.WriteString("</tbody>")
		out.WriteString("</table>")
		if wrap != "" {
			out.WriteString("</div>")
		}
		return out.String()
	}
}
