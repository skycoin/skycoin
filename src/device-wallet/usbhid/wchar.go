// This file is https://github.com/orofarne/gowchar/blob/master/gowchar.go
//
// It was vendored inline to work around CGO limitations that don't allow C types
// to directly cross package API boundaries.
//
// The vendored file is licensed under the 3-clause BSD license, according to:
// https://github.com/orofarne/gowchar/blob/master/LICENSE

// +build !ios
// +build linux darwin windows

package usbhid

/*
#include <wchar.h>

const size_t SIZEOF_WCHAR_T = sizeof(wchar_t);

void gowchar_set (wchar_t *arr, int pos, wchar_t val)
{
	arr[pos] = val;
}

wchar_t gowchar_get (wchar_t *arr, int pos)
{
	return arr[pos];
}
*/
import "C"

import (
	"fmt"
	"unicode/utf16"
	"unicode/utf8"
)

var sizeofWcharT C.size_t = C.size_t(C.SIZEOF_WCHAR_T)

func wcharTToString(s *C.wchar_t) (string, error) {
	switch sizeofWcharT {
	case 2:
		return wchar2ToString(s) // Windows
	case 4:
		return wchar4ToString(s) // Unix
	default:
		panic(fmt.Sprintf("Invalid sizeof(wchar_t) = %v", sizeofWcharT))
	}
}

// Windows
func wchar2ToString(s *C.wchar_t) (string, error) {
	var i int
	var res string
	for {
		ch := C.gowchar_get(s, C.int(i))
		if ch == 0 {
			break
		}
		r := rune(ch)
		i++
		if !utf16.IsSurrogate(r) {
			if !utf8.ValidRune(r) {
				err := fmt.Errorf("Invalid rune at position %v", i)
				return "", err
			}
			res += string(r)
		} else {
			ch2 := C.gowchar_get(s, C.int(i))
			r2 := rune(ch2)
			r12 := utf16.DecodeRune(r, r2)
			if r12 == '\uFFFD' {
				err := fmt.Errorf("Invalid surrogate pair at position %v", i-1)
				return "", err
			}
			res += string(r12)
			i++
		}
	}
	return res, nil
}

// Unix
func wchar4ToString(s *C.wchar_t) (string, error) {
	var i int
	var res string
	for {
		ch := C.gowchar_get(s, C.int(i))
		if ch == 0 {
			break
		}
		r := rune(ch)
		if !utf8.ValidRune(r) {
			err := fmt.Errorf("Invalid rune at position %v", i)
			return "", err
		}
		res += string(r)
		i++
	}
	return res, nil
}
