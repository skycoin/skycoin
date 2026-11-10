package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

func TestWalletRename(t *testing.T) {
	label := "rewards1"
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/csrf", func(w http.ResponseWriter, _ *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"csrf_token": "tok"}))
	})
	mux.HandleFunc("/api/v1/wallet/update", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "tok", r.Header.Get("X-CSRF-Token"))
		require.Equal(t, "test.wlt", r.FormValue("id"))
		label = r.FormValue("label")
		require.NoError(t, json.NewEncoder(w).Encode("success"))
	})
	mux.HandleFunc("/api/v1/wallet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "test.wlt", r.FormValue("id"))
		require.NoError(t, json.NewEncoder(w).Encode(api.WalletResponse{
			Meta: readable.WalletMeta{Filename: "test.wlt", Label: label},
		}))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	saved := apiClient
	apiClient = api.NewClient(srv.URL)
	defer func() { apiClient = saved }()

	require.NoError(t, walletRename("test.wlt", "rewards"))
	require.Equal(t, "rewards", label)
}
