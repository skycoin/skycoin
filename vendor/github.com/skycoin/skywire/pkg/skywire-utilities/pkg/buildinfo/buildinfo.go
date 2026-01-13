// Package buildinfo pkg/skywire-utilities/pkg/buildinfo/buildinfo.go
package buildinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"runtime/debug"
	"strings"
)

const unknown = "unknown"

// Variables set via -ldflags during build
// $ go build -mod=vendor -ldflags="-X 'github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo.version=$(git describe)' -X 'github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo.date=$(date -u "+%Y-%m-%dT%H:%M:%SZ")' -X 'github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo.commit=$(git rev-list -1 HEAD)'" .
var (
	version   = unknown
	commit    = unknown
	date      = unknown
	goversion = ""
)

// format hint: bi.Main.Version = v1.3.29-rc7.0.20250410212328-dc5d22b7ab2a
var bi *debug.BuildInfo

// TODO: deprecate?
// $ go build -ldflags="-X 'github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo.golist=$(go list -m -json -mod=mod github.com/skycoin/<repo>@<branch>)' -X 'github.com/skycoin/skywire/pkg/skywire-utilities/pkg/buildinfo.date=$(date -u "+%Y-%m-%dT%H:%M:%SZ")'" .
// ldflags-provided module info (`go list -m -json`)
var golist string

// ModuleInfo represents the JSON structure returned by `go list -m -json`.
type ModuleInfo struct {
	Version string `json:"Version"`
	Origin  struct {
		Hash string `json:"Hash"`
	} `json:"Origin"`
}

// Regular expressions for commit hash and timestamp
var commitRegex = regexp.MustCompile(`[a-f0-9]{12,}$`) // <-- match commit from end of string
var dateRegex = regexp.MustCompile(`\d{14}`)           // <-- match date anywhere

func init() {
	// Use ldflags-provided `golist` info if available
	if golist != "" {
		var mInfo ModuleInfo
		if err := json.Unmarshal([]byte(golist), &mInfo); err == nil {
			if mInfo.Version != "" && version == unknown {
				version = mInfo.Version
			}
			if mInfo.Origin.Hash != "" && commit == unknown {
				commit = mInfo.Origin.Hash
			}
		}
	}

	// If version is still unknown, try reading from runtime build info
	if version == unknown || version == "" {
		var ok bool
		bi, ok = debug.ReadBuildInfo()
		if ok {
			if bi.Main.Version != "" {
				parseVersionInfo(bi.Main.Version)
			}
			if bi.GoVersion != "" {
				goversion = bi.GoVersion
			}
		}
	}
}

func parseVersionInfo(ver string) {
	// Extract commit
	if match := commitRegex.FindString(ver); match != "" {
		commit = match
		ver = strings.TrimSuffix(ver, "-"+commit)
	}

	// Extract date
	if match := dateRegex.FindString(ver); match != "" {
		date = formatBuildDate(match)
		ver = strings.Replace(ver, match, "", 1)
		ver = strings.TrimSuffix(ver, "-") // remove trailing dash
		ver = strings.TrimSuffix(ver, ".") // remove trailing dot
	}

	// What's left is version
	version = ver
}

// formatBuildDate converts a 14-digit timestamp into RFC3339 format
func formatBuildDate(dateStr string) string {
	if len(dateStr) != 14 {
		return unknown
	}
	return fmt.Sprintf("%s-%s-%sT%s:%s:%sZ",
		dateStr[0:4], dateStr[4:6], dateStr[6:8], // YYYY-MM-DD
		dateStr[8:10], dateStr[10:12], dateStr[12:14], // HH:MM:SS
	)
}

// Version returns the extracted version string.
func Version() string {
	return version
}

// DBIVersion returns bi.Main.Version.
func DBIVersion() string {
	if bi != nil {
		return bi.Main.Version
	}
	return ""
}

// Go returns the Go compiler version used for the build.
func Go() string {
	return goversion
}

// Commit returns the extracted commit hash.
func Commit() string {
	return commit
}

// Date returns the extracted build date in RFC3339 format.
func Date() string {
	return date
}

// DebugBuildInfo returns the raw debug.BuildInfo struct.
func DebugBuildInfo() *debug.BuildInfo {
	return bi
}

// Get returns a summary of build information.
func Get() *Info {
	return &Info{
		Version: Version(),
		Commit:  Commit(),
		Date:    Date(),
		Go:      Go(),
	}
}

// Info represents build metadata.
type Info struct {
	Go      string `json:"go,omitempty"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// WriteTo writes build info summary to an io.Writer.
func (info *Info) WriteTo(w io.Writer) (int64, error) {
	msg := fmt.Sprintf("Version %q built on %q against commit %q\n", info.Version, info.Date, info.Commit)
	n, err := w.Write([]byte(msg))
	return int64(n), err
}
