package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/skycoin/skycoin/cmd/cipher_test_suite/cipherTestSuite"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	inputData *cipherTestSuite.InputData
)

func main() {
	testDataDir := path.Join(os.Getenv("GOPATH"), "/src/github.com/skycoin/skycoin/cmd/cipher_test_suite/testdata/")
	inputData = cipherTestSuite.ReadInputData(path.Join(testDataDir, "inputData.golden"))
	files := traverse(testDataDir, `seed-\d+.golden`)
	for _, f := range files {
		processFile(path.Join(testDataDir, f))
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

func processFile(goldenFile string) {
	var data cipherTestSuite.SeedSignature

	f, err := os.Open(goldenFile)
	if err != nil {
		log.Panicf("failed open file: %v, err: %v", goldenFile, err)
	}
	reader := bufio.NewReader(f)
	defer f.Close()
	err = json.NewDecoder(reader).Decode(&data)
	if err != nil {
		log.Panicf("failed decode seed data from file. err: %v", err)
	}
	for _, dataItem := range data.Keys {
		checkAddress(data.Seed, dataItem.Address)
		for idx, signature := range dataItem.Signatures {
			verifySignature(dataItem.Public, inputData.Hashes[idx], signature)
		}
	}
}

func checkAddress(seed string, address string) {
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		log.Panicf("failed decode seed from string. err: %v", err)
	}
	_, _, addr := cipherTestSuite.GenerateSecPubAddress(seedBytes)
	if addr.String() != address {
		log.Panicf("failed checkAddress. addresses are not equal. need: %v, got: %v", address, addr.String())
	}
}

func verifySignature(publicHex string, hashHex string, signatureHex string) {
	public, err := cipher.PubKeyFromHex(publicHex)
	if err != nil {
		log.Panicf("failed decode public from hex. err: %v", err)
	}
	hash, err := cipher.SHA256FromHex(hashHex)
	if err != nil {
		log.Panicf("failed decode hash from hex. err: %v", err)
	}
	signature, err := cipher.SigFromHex(signatureHex)
	if err != nil {
		log.Panicf("failed decode signature from hex. err: %v", err)
	}
	err = cipher.VerifySignature(public, signature, hash)
	if err != nil {
		log.Panicf("cipher.VerifySignature. publicKey: %v, err: %v", public.Hex(), err)
	}

}
