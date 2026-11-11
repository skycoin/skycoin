package commands

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClassicHandler(t *testing.T) {
	h := classicHandler()

	for _, p := range []string{"/classic/", "/classic/app/blocks/1", "/classic/app/address/2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8/2"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, p, nil))
		body := w.Body.String()
		if w.Code != http.StatusOK || !strings.Contains(body, `<base href="/classic/">`) || !strings.Contains(body, "New explorer</a></body>") {
			t.Errorf("%s: %d, want the index with its base at /classic/ and the link to the new explorer", p, w.Code)
		}
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/classic/favicon.ico", nil))
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "<base") {
		t.Errorf("favicon.ico: %d, want the file", w.Code)
	}
}
