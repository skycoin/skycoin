package readable

import "github.com/blang/semver"

// BuildInfo represents the build info
type BuildInfo struct {
	Version string `json:"version"` // version number
	Commit  string `json:"commit"`  // git commit id
	Branch  string `json:"branch"`  // git branch name
}

// Semver returns the parsed semver.Version of the configured Version string
func (b BuildInfo) Semver() (*semver.Version, error) {
	sv, err := semver.Make(b.Version)
	if err != nil {
		return nil, err
	}

	return &sv, nil
}
