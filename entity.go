package tagfunctions

import (
	"fmt"

	"golang.org/x/net/html"
)

// given a $ent[VALUE] -> &VALUE;
// no error checkign to see if a value HTML or XML entity
func Entity(n *html.Node) error {
	ent := GetArg(n, 0)
	if len(ent) <= 0 || len(ent) > 15 {
		return fmt.Errorf("Got unknown entity %q", ent)
	}
	for _, b := range []byte(ent) {
		switch {
		case b == '#':
		case b >= 'a' && b <= 'z':
		case b >= 'A' && b <= 'Z':
		case b >= '0' && b <= '9':
		default:
			return fmt.Errorf("Got unknown entity %q", ent)
		}
	}
	// ok seems saw, let convert this to a raw node.
	n.Type = html.RawNode
	n.Attr = nil
	n.DataAtom = 0
	n.Data = "&" + ent + ";"
	return nil
}
