package readable

// BuildInfo represents the build info
type BuildInfo struct {
	Version string `json:"version"` // version number
	Commit  string `json:"commit"`  // git commit id
	Branch  string `json:"branch"`  // git branch name
}
