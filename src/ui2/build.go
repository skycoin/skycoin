// Package ui2 src/ui2/build.go c3-ui-wallet
package ui2

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"net/http"
	"runtime/debug"
	"sort"
	"sync"
)

var (
	ownBuildOnce sync.Once
	ownBuild     []byte
)

// ownBuildHash hashes what this package serves: the static files, the
// translations and the commit the templates were compiled from.
func ownBuildHash() []byte {
	ownBuildOnce.Do(func() {
		h := sha256.New()
		for _, fsys := range []fs.FS{staticFS, i18nSub} {
			var names []string
			_ = fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error { //nolint:errcheck
				if err == nil && !d.IsDir() {
					names = append(names, p)
				}
				return nil
			})
			sort.Strings(names)
			for _, n := range names {
				b, _ := fs.ReadFile(fsys, n) //nolint:errcheck
				sum := sha256.Sum256(b)
				h.Write([]byte(n))
				h.Write(sum[:])
			}
		}
		if bi, ok := debug.ReadBuildInfo(); ok {
			for _, s := range bi.Settings {
				if s.Key == "vcs.revision" || s.Key == "vcs.modified" {
					h.Write([]byte(s.Key + "=" + s.Value))
				}
			}
		}
		ownBuild = h.Sum(nil)
	})
	return ownBuild
}

// BuildID names the build an open page runs. It changes when a static file,
// a translation, the commit or the host's own assets (Options.Build) change,
// and ui2.js reloads the page when it does.
func BuildID(o Options) string {
	h := sha256.New()
	h.Write(ownBuildHash())
	if o.Build != nil {
		h.Write([]byte(o.Build()))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// serveBuild answers the poll of ui2.js with the current BuildID.
func serveBuild(o Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(BuildID(o))) //nolint:errcheck
	}
}
