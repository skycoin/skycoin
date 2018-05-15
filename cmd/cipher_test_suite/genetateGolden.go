package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/skycoin/skycoin/cmd/cipher_test_suite/cipherTestSuite"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
)

const (
	hashesCount             = 10
	hashSize                = 32
	seedsCount              = 10
	shortSeedAddressesCount = 10
	longSeedAddressesCount  = 1000
)

var (
	testDataDir       string
	inputData         *cipherTestSuite.InputData
	inputDataFilename string
)

type Job struct {
	JobID          int
	Seed           string
	AddressesCount int
}

func init() {
	testDataDir = path.Join(os.Getenv("GOPATH"),
		"/src/github.com/skycoin/skycoin/cmd/cipher_test_suite/testdata/")
	inputDataFilename = path.Join(testDataDir, "inputData.golden")
}

func main() {
	var job Job
	genInputData := flag.Bool("i", false, "Generate inputData.golden")
	genSeedData := flag.Bool("s", false, "Generate seed-$n.golden")
	flag.Parse()
	if *genInputData {
		generateInputData(inputDataFilename)
	}
	if *genSeedData {
		inputData = cipherTestSuite.ReadInputData(inputDataFilename)
		jobs := make(chan Job, seedsCount)
		results := make(chan bool, seedsCount)
		// generate in parallel to improve speed
		for i := 0; i < runtime.NumCPU(); i++ {
			go worker(jobs, results)
		}
		// generate seed with 1 byte length
		job = Job{
			JobID:          0,
			Seed:           hex.EncodeToString(cipher.RandByte(1)),
			AddressesCount: shortSeedAddressesCount,
		}
		jobs <- job

		for generatedCount := 1; generatedCount < seedsCount; generatedCount++ {
			job := Job{
				JobID:          generatedCount,
				AddressesCount: longSeedAddressesCount,
			}
			// separate seed generation type
			if generatedCount > int(seedsCount/2) {
				seed, err := bip39.NewDefaultMnemomic()
				if err != nil {
					log.Panicf("failed generate seed bip39.NewDefaultMnemomic(). err: %v", err)
				}
				job.Seed = base64.RawStdEncoding.EncodeToString([]byte(seed))
			} else {
				job.Seed = base64.RawStdEncoding.EncodeToString([]byte(cipher.SumSHA256(cipher.RandByte(1024)).Hex()))
			}
			jobs <- job
		}
		close(jobs)
		resultsCount := 0
		for range results {
			resultsCount++
			if resultsCount >= seedsCount {
				break
			}
		}
	}
}

func generateInputData(filename string) {
	inputData := cipherTestSuite.InputData{
		Hashes: make([]string, 0),
	}
	inputData.Hashes = append(inputData.Hashes, cipher.SumSHA256(bytes.Repeat([]byte{0}, hashSize)).Hex())
	inputData.Hashes = append(inputData.Hashes, cipher.SumSHA256(bytes.Repeat([]byte{1}, hashSize)).Hex())
	for i := 0; i < hashesCount-2; i++ {
		inputData.Hashes = append(inputData.Hashes, cipher.SumSHA256(cipher.RandByte(hashSize)).Hex())
	}
	contentJSON, err := json.MarshalIndent(inputData, "", "\t")
	if err != nil {
		log.Panicf("failed encode inputData. err: %v", err)
	}
	err = ioutil.WriteFile(filename, contentJSON, 0644)
	if err != nil {
		log.Panicf("failed to write into file. err: %v", err)
	}
}

func worker(jobs <-chan Job, results chan<- bool) {
	for job := range jobs {
		summary := make(map[string]int)
		data := &cipherTestSuite.SeedSignature{
			Seed: job.Seed,
			Keys: make([]*cipherTestSuite.SeedData, 0),
		}
		log.Printf("job %v/%v\n", job.JobID, seedsCount-1)

		// generate signatures for a part of cases to prevent large .golden files
		for i := 0; i < job.AddressesCount; i++ {
			seedData := generateSeedData([]byte(job.Seed), job.AddressesCount <= 10 || i < int(job.AddressesCount/2))
			summary[seedData.Public]++
			summary[seedData.Secret]++
			summary[seedData.Address]++
			data.Keys = append(data.Keys, seedData)
		}
		// check that all public/secret/address values are equal
		for k, v := range summary {
			if v != job.AddressesCount {
				log.Panicf("generated values are not equal to previous public/secret/address values. key: %v, count: %v", k, v)
			}
		}
		filename := path.Join(testDataDir, fmt.Sprintf("seed-%d.golden", job.JobID))
		contentJSON, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Panicf("failed encode inputData. err: %v", err)
		}
		err = ioutil.WriteFile(filename, contentJSON, 0644)
		if err != nil {
			log.Panicf("failed to write into file. err: %v", err)
		}
		results <- true
	}
}

func generateSeedData(seed []byte, generateSignatures bool) *cipherTestSuite.SeedData {
	secretKey, publicKey, addr := cipherTestSuite.GenerateSecPubAddress(seed)
	data := &cipherTestSuite.SeedData{
		Signatures: make([]string, 0),
		Public:     publicKey.Hex(),
		Secret:     secretKey.Hex(),
		Address:    addr.String(),
	}

	if generateSignatures {
		for _, hash := range inputData.Hashes {
			shaHash, err := cipher.SHA256FromHex(hash)
			if err != nil {
				log.Panicf("failed decode string. err: %v", err)
			}
			data.Signatures = append(data.Signatures, cipher.SignHash(shaHash, secretKey).Hex())
		}
	}
	return data
}
