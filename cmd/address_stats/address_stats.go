/*
address_stats generates statistical data about address generation
*/
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

func zeroPadding(countLen int, value string) string {
	zeros := countLen - len(value)
	for i := 0; i < zeros; i++ {
		value = "0" + value
	}
	return value
}

func main() {
	examples := flag.Int("n", 1, "Number of addresses to generate")
	addrsStats := flag.Bool("addrs", false, "create histogram for address")
	pkeysStats := flag.Bool("pkeys", false, "create histogram for public keys")
	hashesStats := flag.Bool("hashes", false, "create histogram for hashes")
	flag.Parse()

	OneByteMap := make(map[byte]int)
	TwoByteMap := make(map[string]int)

	OneLetterMap := make(map[string]int)
	TwoLetterMap := make(map[string]int)

	OneByteRawMap := make(map[byte]int)
	TwoByteRawMap := make(map[string]int)

	if !*addrsStats && !*pkeysStats && !*hashesStats {
		fmt.Println("you need to choose object for analyze(use one of flags: addrs, pkeys, hashes)")
		return
	}

	start := time.Now()

	pubKeys := make([]cipher.PubKey, *examples)
	rawAddresses := make([]cipher.Ripemd160, *examples)
	addresses := make([]string, *examples)

	//generate pubkeys
	for i := 0; i < *examples; i++ {
		p, _ := cipher.GenerateKeyPair()
		pubKeys[i] = p
	}

	//generate addresses
	for i, p := range pubKeys {
		addresses[i] = cipher.AddressFromPubKey(p).String()
		rawAddresses[i] = cipher.AddressFromPubKey(p).Key
	}

	//analyze addresses
	var output []string
	if *addrsStats {
		for _, a := range addresses {
			if _, ok := OneLetterMap[string([]rune(a)[0])]; ok {
				OneLetterMap[string([]rune(a)[0])]++
			} else {
				OneLetterMap[string([]rune(a)[0])] = 1
			}
			if _, ok := TwoLetterMap[string([]rune(a)[:2])]; ok {
				TwoLetterMap[string([]rune(a)[:2])]++
			} else {
				TwoLetterMap[string([]rune(a)[:2])] = 1
			}
		}

		output = append(output, "\nAddress 1st letter stat:\n")
		for k, v := range OneLetterMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + k + " ]\n"
			output = append(output, s)
		}

		output = append(output, "\nAddress 1-2nd letter stat:\n")
		for k, v := range TwoLetterMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + k + " ]\n"
			output = append(output, s)
		}

	}

	if *pkeysStats {
		for _, p := range pubKeys {
			//first byte gist
			if _, ok := OneByteMap[p[0]]; ok {
				OneByteMap[p[0]]++
			} else {
				OneByteMap[p[0]] = 1
			}
			//2 first byte gist

			if _, ok := TwoByteMap[string(p[:2])]; ok {
				TwoByteMap[string(p[:2])]++
			} else {
				TwoByteMap[string(p[:2])] = 1
			}

		}

		output = append(output, "\nPublic key 1st byte stat:\n")
		for k, v := range OneByteMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			var data []byte
			data = append(data, k)
			bytes := hex.EncodeToString(data)
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + bytes + " ]\n"
			output = append(output, s)
		}

		output = append(output, "\nPublic key 1-2nd byte stat:\n")

		for k, v := range TwoByteMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			bytes := fmt.Sprintf("%x", []byte(k))
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + bytes + " ]\n"
			output = append(output, s)
		}
	}

	if *hashesStats {
		for _, ra := range rawAddresses {
			if _, ok := OneByteRawMap[ra[0]]; ok {
				OneByteRawMap[ra[0]]++
			} else {
				OneByteRawMap[ra[0]] = 1
			}
			//2 first byte gist

			if _, ok := TwoByteRawMap[string(ra[:2])]; ok {
				TwoByteRawMap[string(ra[:2])]++
			} else {
				TwoByteRawMap[string(ra[:2])] = 1
			}

		}

		output = append(output, "\nRaw address 1st byte stat:\n")
		for k, v := range OneByteRawMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			var data []byte
			data = append(data, k)
			bytes := hex.EncodeToString(data)
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + bytes + " ]\n"
			output = append(output, s)
		}

		output = append(output, "\nRaw address 1-2nd byte stat:\n")
		for k, v := range TwoByteRawMap {
			formatV := fmt.Sprintf("%d", v)
			formatV = zeroPadding(len(strconv.Itoa(*examples)), formatV)
			bytes := fmt.Sprintf("%x", []byte(k))
			percent := 100 * float64(v) / float64(*examples)
			per := fmt.Sprintf("%.2f", percent)
			s := formatV + " of " + strconv.Itoa(*examples) + ", " + per + "%, [ " + bytes + " ]\n"
			output = append(output, s)
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)

	f, err := os.Create("histogram")
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range output {
		fmt.Fprint(f, value)
	}

	fmt.Println("Time elapsed: ", elapsed.Seconds())

}
