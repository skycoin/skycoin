package testsuite

import (
	"bufio"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

var (
	inputData *InputData
)

func TestAddressSignature(t *testing.T) {
	testDataDir := path.Join(os.Getenv("GOPATH"), "/src/github.com/skycoin/skycoin/cmd/cipher-testdata/testdata/")
	inputData = ReadInputData(path.Join(testDataDir, "inputData.golden"))
	files := traverse(testDataDir, `seed-\d+.golden`)

	for _, f := range files {
		processFile(t, path.Join(testDataDir, f))
	}
}

func traverse(dir string, filenameTemplate string) []string {
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(filenameTemplate, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func processFile(t *testing.T, goldenFile string) {
	var data SeedSignature
	f, err := os.Open(goldenFile)
	require.NoError(t, err, "failed open file: %v, err: %v", goldenFile, err)
	reader := bufio.NewReader(f)
	defer f.Close()
	err = json.NewDecoder(reader).Decode(&data)
	require.NoError(t, err, "failed decode seed data from file. err: %v", err)
	for _, dataItem := range data.Keys {
		checkAddress(t, data.Seed, dataItem.Address)
		for idx, signature := range dataItem.Signatures {
			verifySignature(t, dataItem.Public, inputData.Hashes[idx], signature)
		}
	}
}

func checkAddress(t *testing.T, seed string, address string) {
	_, _, addr := GenerateSecPubAddress([]byte(seed))
	require.Equal(t, address, addr.String(), "failed checkAddress. addresses are not equal. need: %v, got: %v", address, addr.String())
}

func verifySignature(t *testing.T, publicHex string, hashHex string, signatureHex string) {
	public, err := cipher.PubKeyFromHex(publicHex)
	require.NoError(t, err, "failed decode public from hex. err: %v", err)
	hash, err := cipher.SHA256FromHex(hashHex)
	require.NoError(t, err, "failed decode hash from hex. err: %v", err)
	signature, err := cipher.SigFromHex(signatureHex)
	require.NoError(t, err, "failed decode signature from hex. err: %v", err)
	err = cipher.VerifySignature(public, signature, hash)
	require.NoError(t, err, "cipher.VerifySignature. publicKey: %v, err: %v", public.Hex(), err)
}
