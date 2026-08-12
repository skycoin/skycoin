package buildinfo

import (
	"runtime/debug"
	"testing"
)

// These pin the behavior that keeps regressing: the skycoin commands are
// mounted by other projects, and a version taken from the main module names the
// host rather than skycoin.

func withBuildInfo(t *testing.T, replacement *debug.BuildInfo, fn func()) {
	t.Helper()

	original := bi
	bi = replacement
	defer func() { bi = original }()

	fn()
}

func TestSelfVersionWhenSkycoinIsTheProgram(t *testing.T) {
	withBuildInfo(t, &debug.BuildInfo{
		Main: debug.Module{Path: ModulePath, Version: "v0.28.6"},
	}, func() {
		if got := SelfVersion(); got != "v0.28.6" {
			t.Errorf("SelfVersion() = %q, want the main module's version", got)
		}
	})
}

// `skywire skycoin` runs these commands inside skywire's binary. The banner has
// to name skycoin's version, not skywire's.
func TestSelfVersionWhenEmbeddedInAnotherProgram(t *testing.T) {
	withBuildInfo(t, &debug.BuildInfo{
		Main: debug.Module{Path: "github.com/skycoin/skywire", Version: "v1.3.92"},
		Deps: []*debug.Module{
			{Path: "github.com/some/other", Version: "v1.0.0"},
			{Path: ModulePath, Version: "v0.28.6-0.20260811181324-ab113fbd4466"},
		},
	}, func() {
		want := "v0.28.6-0.20260811181324-ab113fbd4466"
		if got := SelfVersion(); got != want {
			t.Errorf("SelfVersion() = %q, want %q", got, want)
		}

		// The raw accessor keeps meaning the main module, since the --bv flag is
		// documented as printing exactly that.
		if got := DBIVersion(); got != "v1.3.92" {
			t.Errorf("DBIVersion() = %q, want the host's version", got)
		}
	})
}

// A host that does not record us as a dependency leaves nothing better to say.
func TestSelfVersionFallsBackToTheMainModule(t *testing.T) {
	withBuildInfo(t, &debug.BuildInfo{
		Main: debug.Module{Path: "github.com/skycoin/skywire", Version: "v1.3.92"},
	}, func() {
		if got := SelfVersion(); got != "v1.3.92" {
			t.Errorf("SelfVersion() = %q, want the main module's version as a fallback", got)
		}
	})
}

func TestSelfVersionWithoutBuildInfo(t *testing.T) {
	withBuildInfo(t, nil, func() {
		if got := SelfVersion(); got != unknown {
			t.Errorf("SelfVersion() = %q, want %q", got, unknown)
		}
	})
}

// A module built and run straight from a work tree reports "(devel)", which
// names nothing.
func TestSelfVersionRejectsDevel(t *testing.T) {
	withBuildInfo(t, &debug.BuildInfo{
		Main: debug.Module{Path: ModulePath, Version: "(devel)"},
	}, func() {
		if got := SelfVersion(); got != unknown {
			t.Errorf("SelfVersion() = %q, want %q", got, unknown)
		}
	})
}

func TestDepVersion(t *testing.T) {
	withBuildInfo(t, &debug.BuildInfo{
		Main: debug.Module{Path: "github.com/skycoin/skywire", Version: "v1.3.92"},
		Deps: []*debug.Module{{Path: ModulePath, Version: "v0.28.6"}},
	}, func() {
		if got := DepVersion(ModulePath); got != "v0.28.6" {
			t.Errorf("DepVersion(%q) = %q, want v0.28.6", ModulePath, got)
		}

		if got := DepVersion("github.com/not/here"); got != "" {
			t.Errorf("DepVersion of an absent module = %q, want empty", got)
		}
	})
}
