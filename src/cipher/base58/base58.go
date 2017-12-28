// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base58

//Subpackage for encoding data (namely Bitcoin Addresses) into base58 strings

import (
	"encoding/hex"
	"errors"
	"math/big"
)

//Useful materials:
//https://en.bitcoin.it/wiki/Base_58_Encoding
//http://www.strongasanox.co.uk/2011/03/11/base58-encoding-in-python/

//alphabet used by Bitcoins
var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Base58 type to hold the Base58 string
type Base58 string

//reverse alphabet used for quckly converting base58 strings into numbers
var revalp = map[string]int{
	"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7, "9": 8, "A": 9,
	"B": 10, "C": 11, "D": 12, "E": 13, "F": 14, "G": 15, "H": 16, "J": 17, "K": 18, "L": 19,
	"M": 20, "N": 21, "P": 22, "Q": 23, "R": 24, "S": 25, "T": 26, "U": 27, "V": 28, "W": 29,
	"X": 30, "Y": 31, "Z": 32, "a": 33, "b": 34, "c": 35, "d": 36, "e": 37, "f": 38, "g": 39,
	"h": 40, "i": 41, "j": 42, "k": 43, "m": 44, "n": 45, "o": 46, "p": 47, "q": 48, "r": 49,
	"s": 50, "t": 51, "u": 52, "v": 53, "w": 54, "x": 55, "y": 56, "z": 57,
}

// Hex2Big converts hex to big
func Hex2Big(b []byte) *big.Int {
	answer := big.NewInt(0)

	for i := 0; i < len(b); i++ {
		answer.Lsh(answer, 8)
		answer.Add(answer, big.NewInt(int64(b[i])))
	}

	return answer
}

// String2Hex converts string to hex
func String2Hex(s string) []byte {
	answer, _ := hex.DecodeString(s)
	return answer
}

// ToBig convert base58 to big.Int
func (b Base58) ToBig() (*big.Int, error) {
	answer := new(big.Int)
	for i := 0; i < len(b); i++ {
		answer.Mul(answer, big.NewInt(58)) //multiply current value by 58
		c, ok := revalp[string(b[i:i+1])]
		if !ok {
			return nil, errors.New("Invalid base58 character")
		}
		answer.Add(answer, big.NewInt(int64(c))) //add value of the current letter
	}
	return answer, nil
}

// ToInt converts base58 to int
func (b Base58) ToInt() (int, error) {
	answer := 0
	for i := 0; i < len(b); i++ {
		answer *= 58 //multiply current value by 58
		c, ok := revalp[string(b[i:i+1])]
		if !ok {
			return 0, errors.New("Invalid base58 character")
		}
		answer += c //add value of the current letter
	}
	return answer, nil
}

//ToHex converts base58 to hex bytes
func (b Base58) ToHex() ([]byte, error) {
	value, err := b.ToBig() //convert to big.Int
	if err != nil {
		return nil, err
	}
	oneCount := 0
	bs := string(b)
	if len(bs) == 0 {
		return nil, errors.New("Invalid base58 string")
	}
	for bs[oneCount] == '1' {
		oneCount++
		if oneCount >= len(bs) {
			return nil, errors.New("Invalid base58 string")
		}
	}
	//convert big.Int to bytes
	return append(make([]byte, oneCount), value.Bytes()...), nil
}

// Base582Big converts base58 to big
func (b Base58) Base582Big() (*big.Int, error) {
	answer := new(big.Int)
	for i := 0; i < len(b); i++ {
		answer.Mul(answer, big.NewInt(58)) //multiply current value by 58
		c, ok := revalp[string(b[i:i+1])]
		if !ok {
			return nil, errors.New("Invalid base58 character")
		}
		answer.Add(answer, big.NewInt(int64(c))) //add value of the current letter
	}
	return answer, nil
}

// Base582Int converts base58 to int
func (b Base58) Base582Int() (int, error) {
	answer := 0
	for i := 0; i < len(b); i++ {
		answer *= 58 //multiply current value by 58
		c, ok := revalp[string(b[i:i+1])]
		if !ok {
			return 0, errors.New("Invalid base58 character")
		}
		answer += c //add value of the current letter
	}
	return answer, nil
}

// Base582Hex converts base58 to hex bytes
func Base582Hex(b string) ([]byte, error) {
	return Base58(b).ToHex()
}

// BitHex converts base58 to hexes used by Bitcoins (keeping the zeroes on the front, 25 bytes long)
func (b Base58) BitHex() ([]byte, error) {
	value, err := b.ToBig() //convert to big.Int
	if err != nil {
		return nil, err
	}

	tmp := value.Bytes() //convert to hex bytes
	if len(tmp) == 25 {  //if it is exactly 25 bytes, return
		return tmp, nil
	} else if len(tmp) > 25 { //if it is longer than 25, return nothing
		return nil, errors.New("base58 invalid length")
	}
	answer := make([]byte, 25)      //make 25 byte container
	for i := 0; i < len(tmp); i++ { //copy converted bytes
		answer[24-i] = tmp[len(tmp)-1-i]
	}
	return answer, nil
}

// Big2Base58 encodes big.Int to base58 string
func Big2Base58(val *big.Int) Base58 {
	answer := ""
	valCopy := new(big.Int).Abs(val) //copies big.Int

	if val.Cmp(big.NewInt(0)) <= 0 { //if it is less than 0, returns empty string
		return Base58("")
	}

	tmpStr := ""
	tmp := new(big.Int)
	for valCopy.Cmp(big.NewInt(0)) > 0 { //converts the number into base58
		tmp.Mod(valCopy, big.NewInt(58))                //takes modulo 58 value
		valCopy.Div(valCopy, big.NewInt(58))            //divides the rest by 58
		tmpStr += alphabet[tmp.Int64() : tmp.Int64()+1] //encodes
	}
	for i := (len(tmpStr) - 1); i > -1; i-- {
		answer += tmpStr[i : i+1] //reverses the order
	}
	return Base58(answer) //returns
}

// Int2Base58 encodes int to base58 string
func Int2Base58(val int) Base58 {
	answer := ""

	if val <= 0 { //if it is less than 0, returns empty string
		return Base58("")
	}
	valCopy := val

	tmpStr := ""
	tmp := 0
	for valCopy > 0 { //converts the number into base58
		tmp = valCopy % 58 //takes modulo 58 value
		valCopy /= 58      //divides the rest by 58

		tmpStr += alphabet[tmp : tmp+1] //encodes
	}
	for i := (len(tmpStr) - 1); i > -1; i-- {
		answer += tmpStr[i : i+1] //reverses the order
	}

	return Base58(answer) //returns
}

// Hex2Base58 encodes hex bytes into base58
func Hex2Base58(val []byte) Base58 {
	tmp := Big2Base58(Hex2Big(val)) //encoding of the number without zeroes in front

	//looking for zeros at the beginning
	i := 0
	for i = 0; val[i] == 0 && i < len(val); i++ {
	}
	answer := ""
	for j := 0; j < i; j++ { //adds zeroes from the front
		answer += alphabet[0:1]
	}
	answer += string(tmp) //concatenates

	return Base58(answer) //returns
}

// Hex2Base58String converts hex to Base58 string
func Hex2Base58String(val []byte) string {
	return string(Hex2Base58(val))
}

// Hex2Base58Str converts hex to Base58 string
func Hex2Base58Str(val []byte) string {
	return string(Hex2Base58(val))
}
