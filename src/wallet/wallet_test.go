package wallet

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/stretchr/testify/require"
)

type fakeWalletForGuardView struct {
	*MockWallet
	seed      string
	label     string
	n         int
	encrypted bool
}

func (f fakeWalletForGuardView) Label() string {
	return f.label
}

func (f fakeWalletForGuardView) Seed() string {
	return f.seed
}

func (f fakeWalletForGuardView) IsEncrypted() bool {
	return f.encrypted
}

func (f fakeWalletForGuardView) Unlock(pwd []byte) (Wallet, error) {
	nf := f
	nf.encrypted = false
	return &nf, nil
}

func (f *fakeWalletForGuardView) Erase() {
	f.seed = ""
}

func TestWalletGuard(t *testing.T) {
	tt := []struct {
		name      string
		encrypted bool
		pwd       []byte
		err       error
	}{
		{
			name:      "ok",
			encrypted: true,
			pwd:       []byte("pwd"),
		},
		{
			name:      "wallet is not encrypted",
			encrypted: false,
			err:       ErrWalletNotEncrypted,
		},
		{
			name:      "password is nil",
			encrypted: true,
			pwd:       []byte(""),
			err:       ErrMissingPassword,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			seed := bip39.MustNewDefaultMnemonic()
			w := &fakeWalletForGuardView{
				seed:      seed,
				label:     "label",
				encrypted: tc.encrypted,
			}

			err := GuardView(w, tc.pwd, func(wlt Wallet) error {
				require.False(t, wlt.IsEncrypted())
				return nil
			})

			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			require.True(t, w.IsEncrypted())
		})
	}
}

func TestRemoveBackupFiles(t *testing.T) {
	type wltInfo struct {
		wltName string
		version string
	}

	tt := []struct {
		name                   string
		initFiles              []wltInfo
		expectedRemainingFiles map[string]struct{}
	}{
		{
			name:                   "no file",
			initFiles:              []wltInfo{},
			expectedRemainingFiles: map[string]struct{}{},
		},
		{
			name: "wlt v0.1=1 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=2 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=2 delete 2 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=3 delete 3 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=1 no delete",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t4.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t4.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.2=3 bak v0.2=1 no delete",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.2",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.1=1 bak v0.1=1 wlt v0.2=2 bak v0.2=2 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t2.wlt.bak",
					"0.2",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t2.wlt.bak": struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.1=1 bak v0.1=2 wlt v0.2=2 bak v0.2=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t2.wlt.bak": struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			// Initialize files
			mockWltContentTmp := `
{
    "meta": {
        "type": "deterministic",
        "version": "{{.Version}}"
    },
    "entries": []
}`

			tmp := template.New("mockWalletCntTmp")
			tmp, err := tmp.Parse(mockWltContentTmp)
			require.NoError(t, err)

			for _, f := range tc.initFiles {
				fw, err := os.Create(filepath.Join(dir, f.wltName))
				defer fw.Close()
				err = tmp.Execute(fw, struct{ Version string }{f.version})
				require.NoError(t, err)
			}

			require.NoError(t, removeBackupFiles(dir))

			// Get all remaining files
			fs, err := ioutil.ReadDir(dir)
			require.NoError(t, err)
			require.Len(t, fs, len(tc.expectedRemainingFiles))
			for _, f := range fs {
				_, ok := tc.expectedRemainingFiles[f.Name()]
				require.True(t, ok)
			}
		})
	}
}

func prepareWltDir() string {
	dir, err := ioutil.TempDir("", "wallets")
	if err != nil {
		panic(err)
	}

	return dir
}
