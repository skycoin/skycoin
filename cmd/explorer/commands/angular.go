package commands

import (
	"bytes"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/skycoin/skycoin/explorer"
)

// angularPrefix is where the Angular explorer is served
// beside the server rendered one.
const angularPrefix = "/angular/"

// toNew is a link from the Angular explorer to the same page in the server
// rendered one. The href is filled when clicked, after the Angular router
// has moved.
const toNew = `<a href="/" onclick="this.href=location.pathname.replace(/^\/angular/,'')" ` +
	`style="position:fixed;right:12px;bottom:12px;z-index:1000;padding:6px 12px;` +
	`border-radius:16px;background:#0072ff;color:#fff;font:14px sans-serif;text-decoration:none">` +
	`New explorer</a>`

// angular serves the Angular explorer below angularPrefix: its files where
// they exist, else its index.html with <base> moved there, so its router
// sees its usual /app paths.
type angular struct {
	fsys  fs.FS
	files http.Handler
}

func angularHandler() angular {
	var fsys fs.FS
	if uiFilesFolder != "" {
		fsys = os.DirFS(uiFilesFolder)
	} else if sub, err := fs.Sub(explorer.EmbeddedFiles, "dist"); err == nil {
		fsys = sub
	} else {
		fsys = explorer.EmbeddedFiles
	}
	return angular{fsys: fsys, files: http.FileServer(http.FS(fsys))}
}

func (c angular) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(path.Clean(r.URL.Path), strings.TrimSuffix(angularPrefix, "/"))
	if st, err := fs.Stat(c.fsys, strings.TrimPrefix(name, "/")); err == nil && !st.IsDir() {
		http.StripPrefix(strings.TrimSuffix(angularPrefix, "/"), c.files).ServeHTTP(w, r)
		return
	}
	b, err := fs.ReadFile(c.fsys, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	b = bytes.Replace(b, []byte(`<base href="/">`), []byte(`<base href="`+angularPrefix+`">`), 1)
	b = bytes.Replace(b, []byte("</body>"), []byte(toNew+"</body>"), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(b); err != nil {
		log.Printf("Error serving index.html: %v", err)
	}
}
