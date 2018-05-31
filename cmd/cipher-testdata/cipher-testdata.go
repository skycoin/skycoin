package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/cipher/testsuite"
	"github.com/skycoin/skycoin/src/util/file"
)

const (
	inputTestDataFilename = "input-hashes.golden"
	manyAddressesFilename = "many-addresses.golden"
	seedFilenameFormat    = "seed-%04d.golden"
	randomSeedLength      = 1024
)

var help = fmt.Sprintf(`cipher-testdata generates testdata to be used by the cipher test suite in src/cipher/testsuite.

A file named %s will be generated,
which contains a list of hex-encoded hashes to sign.
This list of hashes will always include a hash whose bytes are all 0x00,
and a hash whose bytes are all 0xFF.

Multiple files named seed-{num}.json will be generated.
Each of these files contains a seed, a number of secret keys,
public keys and addresses generated from this seed.
For each secret key, each hash from inputs will be signed,
and the result saved to the file.
Half of the seeds will be generated as SHA256(RandByte(1024)) and half will
be generated as bip39 seeds. Seeds are base64 encoded in the JSON file.

A seed of length 1 is always generated,
in addition to the requested number of seeds.

A file named %s will be generated,
which contains a seed and a number of secret keys,
public keys and addresses generated from this seed.
The number of secret keys generated is much larger than for the other seeds.
This file is used to test deterministic key generation more thoroughly.
This file will not contain any signatures,
because the filesize would be too large.`, inputTestDataFilename, manyAddressesFilename)

type job struct {
	jobID        int
	seed         []byte
	addressCount int
}

func init() {
	flag.Usage = func() {
		// TODO go1.10 - use flag.CommandLine.Output() (not support in go1.9)
		// fmt.Fprintf(flag.CommandLine.Output(), "%s\n\nUsage of %s:\n", help, os.Args[0])
		fmt.Fprintf(os.Stderr, "%s\n\nUsage of %s:\n", help, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var j job

	seedsCount := flag.Int("seeds", 10, "number of seeds to generate")
	inputsCount := flag.Int("hashes", 8, "number of random hashes for input-hashes.golden")
	addressCount := flag.Int("addresses", 10, "number of addresses to generate per seed")
	manyAddressesCount := flag.Int("many-addresses", 1000, "number of addresses to generate for the single many-addresses test data")
	outputDir := flag.String("dir", "./testdata", "output directory")

	flag.Parse()

	fmt.Println("Creating output directory", *outputDir)

	// Create the output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Generating", manyAddressesFilename)

	// Generate the many-addresses testdata
	manyAddressesData := generateSeedTestData(job{
		seed:         []byte(bip39.MustNewDefaultMnemonic()),
		addressCount: *manyAddressesCount,
	})
	fn := filepath.Join(*outputDir, manyAddressesFilename)
	if err := file.SaveJSON(fn, manyAddressesData.ToJSON(), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Generating", inputTestDataFilename)

	// Create the input hashes used for signing
	inputs := generateInputTestData(*inputsCount)
	fn = filepath.Join(*outputDir, inputTestDataFilename)
	if err := file.SaveJSON(fn, inputs.ToJSON(), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Generating seed data times", *seedsCount)

	jobs := make([]job, 0, *seedsCount+1)

	// Generate seed with 1 byte length
	jobs = append(jobs, job{
		seed:         cipher.RandByte(1),
		addressCount: *addressCount,
	})

	// Generate random and mnemonic seeds
	for i := 0; i < *seedsCount; i++ {
		j = job{
			addressCount: *addressCount,
		}

		if i%2 == 0 {
			j.seed = []byte(bip39.MustNewDefaultMnemonic())
		} else {
			hash := cipher.SumSHA256(cipher.RandByte(randomSeedLength))
			j.seed = hash[:]
		}

		jobs = append(jobs, j)
	}

	seedTestData := make(chan *testsuite.SeedTestData, len(jobs))
	writeDone := make(chan struct{})

	go func() {
		defer close(writeDone)

		var i int
		for data := range seedTestData {
			filename := filepath.Join(*outputDir, fmt.Sprintf(seedFilenameFormat, i))
			if err := file.SaveJSON(filename, data.ToJSON(), 0644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			i++
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for i, j := range jobs {
		j.jobID = i
		go func() {
			defer wg.Done()
			data := generateSeedTestData(j)
			signSeedTestData(data, inputs.Hashes)
			seedTestData <- data
		}()
	}
	wg.Wait()

	close(seedTestData)

	<-writeDone
}

func generateInputTestData(inputsCount int) *testsuite.InputTestData {
	var hashes []cipher.SHA256

	// Add a hash which is all zeroes
	hashes = append(hashes, cipher.SumSHA256(bytes.Repeat([]byte{0}, 32)))
	// Add a hash which is all ones
	hashes = append(hashes, cipher.SumSHA256(bytes.Repeat([]byte{1}, 32)))

	for i := 0; i < inputsCount; i++ {
		hashes = append(hashes, cipher.SumSHA256(cipher.RandByte(32)))
	}

	return &testsuite.InputTestData{
		Hashes: hashes,
	}
}

func generateSeedTestData(j job) *testsuite.SeedTestData {
	data := &testsuite.SeedTestData{
		Seed: j.seed,
		Keys: make([]testsuite.KeysTestData, j.addressCount),
	}

	keys := cipher.GenerateDeterministicKeyPairs(j.seed, j.addressCount)

	for i, s := range keys {
		data.Keys[i].Secret = s

		p := cipher.PubKeyFromSecKey(s)
		data.Keys[i].Public = p

		addr := cipher.AddressFromPubKey(p)
		data.Keys[i].Address = addr
	}

	return data
}

func signSeedTestData(data *testsuite.SeedTestData, hashes []cipher.SHA256) {
	for i := range data.Keys {
		for _, h := range hashes {
			sig := cipher.SignHash(h, data.Keys[i].Secret)
			data.Keys[i].Signatures = append(data.Keys[i].Signatures, sig)
		}
	}
}
