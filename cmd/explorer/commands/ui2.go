package commands

import (
	"bytes"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin/explorer"
)

// toUI2 is a link from the Angular explorer to the same page in /ui2. The
// href is filled when clicked, after the Angular router has moved.
const toUI2 = `<a href="/ui2/" onclick="this.href='/ui2'+location.pathname" ` +
	`style="position:fixed;right:12px;bottom:12px;z-index:1000;padding:6px 12px;` +
	`border-radius:16px;background:#0072ff;color:#fff;font:14px sans-serif;text-decoration:none">` +
	`New explorer</a>`

// serveClassicIndex serves the Angular index.html with the link to /ui2.
func serveClassicIndex(w http.ResponseWriter) {
	var b []byte
	var err error
	if uiFilesFolder != "" {
		b, err = os.ReadFile(filepath.Join(uiFilesFolder, "index.html")) //nolint:gosec
	} else {
		b, err = fs.ReadFile(explorer.EmbeddedFiles, "dist/index.html")
	}
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	b = bytes.Replace(b, []byte("</body>"), []byte(toUI2+"</body>"), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(b); err != nil {
		log.Printf("Error serving index.html: %v", err)
	}
}
