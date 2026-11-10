package cli

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
)

// nodeStub answers the two requests nodeWallet makes, as a node would.
func nodeStub(t *testing.T, folder func(http.ResponseWriter), wlt func(http.ResponseWriter)) *api.Client {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/wallets/folderName", func(w http.ResponseWriter, _ *http.Request) { folder(w) })
	mux.HandleFunc("/api/v1/wallet", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "test.wlt", r.FormValue("id"))
		wlt(w)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return api.NewClient(srv.URL)
}

func TestNodeWallet(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.wlt")
	sendJSON := func(v interface{}) func(http.ResponseWriter) {
		return func(w http.ResponseWriter) { require.NoError(t, json.NewEncoder(w).Encode(v)) }
	}
	inDir := sendJSON(api.WalletFolder{Address: dir})
	held := sendJSON(api.WalletResponse{})
	never := func(http.ResponseWriter) { t.Fatal("the node was asked for a wallet outside its folder") }

	closed := httptest.NewServer(http.NotFoundHandler())
	closed.Close()

	cases := []struct {
		name   string
		c      interface{}
		file   string
		wantID string
		err    bool
	}{
		{"the node holds the wallet", nodeStub(t, inDir, held), file, "test.wlt", false},
		{"a relative path to the node's folder", nodeStub(t, inDir, held), filepath.Join(dir, "sub", "..", "test.wlt"), "test.wlt", false},
		{"no node is reachable", api.NewClient(closed.URL), file, "", false},
		{"the file is outside the node's folder", nodeStub(t, sendJSON(api.WalletFolder{Address: filepath.Join(dir, "other")}), never), file, "", false},
		{"the WALLET API set is off", nodeStub(t, func(w http.ResponseWriter) { wh.Error403(w, "Endpoint is disabled") }, never), file, "", false},
		{"the node did not load the file", nodeStub(t, inDir, func(w http.ResponseWriter) { wh.Error400(w, wallet.ErrWalletNotExist.Error()) }), file, "", false},
		{"the folder request fails", nodeStub(t, func(w http.ResponseWriter) { wh.Error401(w, "", "Unauthorized") }, never), file, "", true},
		{"the wallet request fails", nodeStub(t, inDir, func(w http.ResponseWriter) { wh.Error500(w, "boom") }), file, "", true},
		{"a client that is not a node", fakeOutputs{}, file, "", false},
		{"no client", nil, file, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n, id, err := nodeWallet(tc.c, tc.file)
			if tc.err {
				require.ErrorContains(t, err, "cannot tell whether the node holds wallet")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantID, id)
			require.Equal(t, tc.wantID != "", n != nil)
		})
	}
}

// fakeNode holds the wallet file of every path under dir, and adds addresses to
// its own copy like a node whose wallet has one address.
type fakeNode struct {
	fakeOutputs
	dir       string
	added     []string
	passwords []string
	change    []string
	err       error
}

func (f *fakeNode) WalletFolderName() (*api.WalletFolder, error) {
	return &api.WalletFolder{Address: f.dir}, nil
}

func (f *fakeNode) Wallet(_ string) (*api.WalletResponse, error) {
	resp := &api.WalletResponse{Accounts: []readable.WalletAccount{{}}}
	for _, a := range f.change {
		resp.Accounts[0].ChangeEntries = append(resp.Accounts[0].ChangeEntries, readable.WalletEntry{Address: a})
	}
	return resp, nil
}

func (f *fakeNode) NewWalletAddress(_ string, password string, options ...wallet.Option) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	var opts wallet.Bip44EntriesOptions
	for _, o := range options {
		o(&opts)
	}
	addr := testutil.MakeAddress().String()
	f.added = append(f.added, addr)
	f.passwords = append(f.passwords, password)
	if opts.ChainMode == wallet.ChangeChain {
		f.change = append(f.change, addr)
	}
	return []string{addr}, nil
}

func TestChangeAddressFromNode(t *testing.T) {
	const seed = "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	to := []SendAmount{{Addr: testutil.MakeAddress().String(), Coins: 1e6}}

	newFile := func(t *testing.T, typ string, password []byte) (string, cipher.Address, []byte) {
		dir := t.TempDir()
		w, err := wallet.NewWallet("test.wlt", "label", seed, wallet.Options{
			Type:       typ,
			GenerateN:  1,
			Encrypt:    len(password) != 0,
			Password:   password,
			CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure,
		})
		require.NoError(t, err)
		require.NoError(t, wallet.Save(w, dir))
		first, err := w.GetEntryAt(0)
		require.NoError(t, err)
		file := filepath.Join(dir, "test.wlt")
		before, err := os.ReadFile(file) //nolint:gosec
		require.NoError(t, err)
		return file, first.Address.(cipher.Address), before
	}
	unchanged := func(t *testing.T, file string, before []byte) {
		after, err := os.ReadFile(file) //nolint:gosec
		require.NoError(t, err)
		require.Equal(t, string(before), string(after), "the node saves the wallet file, the CLI must not")
	}

	cases := []struct {
		name     string
		typ      string
		password []byte
		sent     string
	}{
		{"deterministic", wallet.WalletTypeDeterministic, nil, ""},
		{"encrypted deterministic", wallet.WalletTypeDeterministic, []byte("pwd"), "pwd"},
		{"bip44", wallet.WalletTypeBip44, nil, ""},
		{"encrypted bip44", wallet.WalletTypeBip44, []byte("pwd"), ""},
	}
	for _, tc := range cases {
		t.Run("send "+tc.name, func(t *testing.T) {
			file, first, before := newFile(t, tc.typ, tc.password)
			node := &fakeNode{fakeOutputs: fakeOutputs{t, first}, dir: filepath.Dir(file)}

			txn, err := CreateRawTxnFromWallet(node, file, "", to, NewPasswordReader(tc.password), params.MainNetDistribution)
			require.NoError(t, err)
			require.Len(t, txn.Out, 2)
			require.Len(t, node.added, 1)
			require.Equal(t, node.added[0], txn.Out[1].Address.String(), "the change goes to the address the node added")
			require.Equal(t, []string{tc.sent}, node.passwords)
			require.Equal(t, tc.typ == wallet.WalletTypeBip44, len(node.change) == 1)
			unchanged(t, file, before)
		})
	}

	t.Run("sweep", func(t *testing.T) {
		file, first, before := newFile(t, wallet.WalletTypeDeterministic, nil)
		node := &fakeNode{fakeOutputs: fakeOutputs{t, first}, dir: filepath.Dir(file)}

		txn, dest, err := sweepWallet(node, file, "", nil, params.MainNetDistribution)
		require.NoError(t, err)
		require.Equal(t, node.added, []string{dest})
		require.Equal(t, dest, txn.Out[0].Address.String())
		unchanged(t, file, before)
	})

	t.Run("the node refuses", func(t *testing.T) {
		file, first, before := newFile(t, wallet.WalletTypeDeterministic, nil)
		node := &fakeNode{fakeOutputs: fakeOutputs{t, first}, dir: filepath.Dir(file), err: errors.New("403 Forbidden")}

		_, err := CreateRawTxnFromWallet(node, file, "", to, nil, params.MainNetDistribution)
		require.ErrorContains(t, err, "the node holds wallet test.wlt and could not add an address to it")
		_, _, err = sweepWallet(node, file, "", nil, params.MainNetDistribution)
		require.ErrorContains(t, err, "the node holds wallet test.wlt and could not add an address to it")
		unchanged(t, file, before)
	})

	t.Run("the file is elsewhere", func(t *testing.T) {
		file, first, before := newFile(t, wallet.WalletTypeDeterministic, nil)
		node := &fakeNode{fakeOutputs: fakeOutputs{t, first}, dir: t.TempDir()}

		_, err := CreateRawTxnFromWallet(node, file, "", to, nil, params.MainNetDistribution)
		require.NoError(t, err)
		require.Empty(t, node.added)
		after, err := os.ReadFile(file) //nolint:gosec
		require.NoError(t, err)
		require.NotEqual(t, string(before), string(after), "the new address is saved in the file")
	})
}

func TestGetChangeAddressBip44FromNode(t *testing.T) {
	dir := t.TempDir()
	w, err := bip44wallet.NewWallet("b.wlt", "test", "voyage say extend find sheriff surge priority merit ignore maple cash argue", "")
	require.NoError(t, err)
	require.NoError(t, wallet.Save(w, dir))
	path := filepath.Join(dir, "b.wlt")
	before, err := os.ReadFile(path) //nolint:gosec
	require.NoError(t, err)
	wa := walletAddress{Wallet: path}
	node := &fakeNode{dir: dir}

	first, err := getChangeAddress(node, wa, "", fakeFinder{})
	require.NoError(t, err)
	require.Equal(t, node.change, []string{first}, "the node adds the first change address")

	again, err := getChangeAddress(node, wa, "", fakeFinder{})
	require.NoError(t, err)
	require.Equal(t, first, again, "an unused change address of the node's copy is reused")

	next, err := getChangeAddress(node, wa, "", fakeFinder{first: true})
	require.NoError(t, err)
	require.Equal(t, node.change, []string{first, next}, "a used one moves on")

	after, err := os.ReadFile(path) //nolint:gosec
	require.NoError(t, err)
	require.Equal(t, string(before), string(after))
}
