package readable

import "github.com/blang/semver"

// BuildInfo represents the build info
// swagger:model buildInfo
type BuildInfo struct {
	// version number
	Version string `json:"version"`
	// git commit id
	Commit  string `json:"commit"`
	// git branch name
	Branch  string `json:"branch"`
}

// Semver returns the parsed semver.Version of the configured Version string
func (b BuildInfo) Semver() (*semver.Version, error) {
	sv, err := semver.Make(b.Version)
	if err != nil {
		return nil, err
	}

	return &sv, nil
}
