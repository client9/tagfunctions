package tagfunctions

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// may or may not be needed for CORS + Fonts
func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".woff2") {
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
		}
		fs.ServeHTTP(w, r)
	}
}

func Serve(outdir string, prefix string) {
	fmt.Println("Starting")

	if prefix == "" {
		http.FileServer(http.Dir(outdir))
	} else {
		prefix = "/" + prefix + "/"
		http.Handle(prefix, addHeaders(http.StripPrefix(prefix, http.FileServer(http.Dir(outdir)))))
	}

	if err := http.ListenAndServe("localhost:1313", nil); err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}
}
