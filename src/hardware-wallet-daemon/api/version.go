package api

import (
	"net/http"

	"github.com/blang/semver"
)

// BuildInfo represents the build info
type BuildInfo struct {
	Version string `json:"version"` // version number
	Commit  string `json:"commit"`  // git commit id
	Branch  string `json:"branch"`  // git branch name
}

// versionHandler returns app version data
// URI: /api/v1/version
// Method: GET
func versionHandler(c muxConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: c.build,
		})
	}
}

// Semver returns the parsed semver.Version of the configured Version string
func (b BuildInfo) Semver() (*semver.Version, error) {
	sv, err := semver.Make(b.Version)
	if err != nil {
		return nil, err
	}

	return &sv, nil
}
