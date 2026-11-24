package ui2

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildEndpoint(t *testing.T) {
	wasm := "a"
	o := Options{Build: func() string { return wasm }}
	h := Handler(testNode(), o)

	rr := get(t, h, "/build")
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "no-store", rr.Header().Get("Cache-Control"))
	id := rr.Body.String()
	require.Regexp(t, `^[0-9a-f]{16}$`, id)

	page := get(t, h, "/status").Body.String()
	m := regexp.MustCompile(`data-build="([0-9a-f]+)"`).FindStringSubmatch(page)
	require.NotNil(t, m)
	require.Equal(t, id, m[1])
	require.Contains(t, page, `id="update-bar"`)

	// A new cipher wasm is a new build, so open pages reload.
	wasm = "b"
	require.NotEqual(t, id, get(t, h, "/build").Body.String())
	wasm = "a"
	require.Equal(t, id, get(t, h, "/build").Body.String())

	require.NotEqual(t, id, BuildID(Options{}))

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/build", nil))
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestFlowsHoldReload(t *testing.T) {
	h := Handler(testNode(), Options{})
	require.Contains(t, get(t, h, "/new").Body.String(), "data-no-reload")
	require.Contains(t, get(t, h, "/bw/new").Body.String(), "data-no-reload")
	require.NotContains(t, get(t, h, "/status").Body.String(), "data-no-reload")
}
