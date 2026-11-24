// Copyright 2020 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !(linux && (amd64 || arm64 || loong64 || ppc64le || s390x || riscv64 || 386 || arm))

package libc // import "modernc.org/libc"

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

const (
	modNone = iota
	modHH
	modH
	modL
	modLL
	modLD
	modQ
	modCapitalL
	modJ
	modZ
	modCapitalZ
	modT
	mod32
	mod64
)

// Format of the format string
//
// The format string is a character string, beginning and ending in its initial
// shift state, if any.  The format string is composed of zero or more
// directives: ordinary  characters  (not  %), which  are  copied unchanged to
// the output stream; and conversion specifications, each of which results in
// fetching zero or more subsequent arguments.
func printf(format, args uintptr) []byte {
	// format0 := format
	// args0 := args
	buf := bytes.NewBuffer(nil)
	for {
		switch c := *(*byte)(unsafe.Pointer(format)); c {
		case '%':
			format = printfConversion(buf, format, &args)
		case 0:
			// 			if dmesgs {
			// 				dmesg("%v: %q, %#x -> %q", origin(1), GoString(format0), args0, buf.Bytes())
			// 			}
			return buf.Bytes()
		default:
			format++
			buf.WriteByte(c)
		}
	}
}

// Each conversion specification is introduced by the character %, and ends
// with a conversion specifier.  In between there may be (in this order) zero
// or more flags, an optional minimum field width, an optional  precision  and
// an optional length modifier.
func printfConversion(buf *bytes.Buffer, format uintptr, args *uintptr) uintptr {
	format++ // '%'
	spec := "%"

	// Flags characters
	//
	// The character % is followed by zero or more of the following flags:
flags:
	for {
		switch c := *(*byte)(unsafe.Pointer(format)); c {
		case '#':
			// The value should be converted to an "alternate form".  For o conversions,
			// the first character of the output string is made zero (by prefixing a 0 if
			// it was not zero already).  For x and  X  conversions,  a nonzero result has
			// the string "0x" (or "0X" for X conversions) prepended to it.  For a, A, e,
			// E, f, F, g, and G conversions, the result will always contain a decimal
			// point, even if no digits follow it (normally, a decimal point appears in the
			// results of those conversions only if a digit follows).  For g and G
			// conversions, trailing  zeros are not removed from the result as they would
			// otherwise be.  For other conversions, the result is undefined.
			format++
			spec += "#"
		case '0':
			// The  value  should  be zero padded.  For d, i, o, u, x, X, a, A, e, E, f, F,
			// g, and G conversions, the converted value is padded on the left with zeros
			// rather than blanks.  If the 0 and - flags both appear, the 0 flag is
			// ignored.  If a precision is given with a numeric conversion (d, i, o, u, x,
			// and X), the 0 flag is ignored.  For other conversions, the  behav‐ ior is
			// undefined.
			format++
			spec += "0"
		case '-':
			// The  converted value is to be left adjusted on the field boundary.  (The
			// default is right justification.)  The converted value is padded on the right
			// with blanks, rather than on the left with blanks or zeros.  A - overrides a
			// 0 if both are given.
			format++
			spec += "-"
		case ' ':
			// A blank should be left before a positive number (or empty string) produced
			// by a signed conversion.
			format++
			spec += " "
		case '+':
			// A sign (+ or -) should always be placed before a number produced by a signed
			// conversion.  By default, a sign is used only for negative numbers.  A +
			// overrides a space  if  both  are used.
			format++
			spec += "+"
		default:
			break flags
		}
	}
	flags := spec[1:]
	format, width, hasWidth := parseFieldWidth(format, args)
	if hasWidth {
		spec += strconv.Itoa(width)
	}
	format, prec, hasPrecision := parsePrecision(format, args)
	format, mod := parseLengthModifier(format)

	var str string

more:
	// Conversion specifiers
	//
	// A character that specifies the type of conversion to be applied.  The
	// conversion specifiers and their meanings are:
	switch c := *(*byte)(unsafe.Pointer(format)); c {
	case 'd', 'i':
		// The  int argument is converted to signed decimal notation.  The precision,
		// if any, gives the minimum number of digits that must appear; if the
		// converted value requires fewer digits, it is padded on the left with zeros.
		// The default precision is 1.  When 0 is printed with an explicit precision 0,
		// the output is empty.
		format++
		var arg int64
		if isWindows && mod == modL {
			mod = modNone
		}
		switch mod {
		case modL, modLL, mod64, modJ:
			arg = VaInt64(args)
		case modH:
			arg = int64(int16(VaInt32(args)))
		case modHH:
			arg = int64(int8(VaInt32(args)))
		case mod32, modNone:
			arg = int64(VaInt32(args))
		case modT:
			arg = int64(VaInt64(args))
		default:
			panic(todo("", mod))
		}

		u := uint64(arg)
		if arg < 0 {
			u = -u
		}
		str = printfInteger(arg < 0, u, c, flags, width, intPrec(hasPrecision, prec))
	case 'u':
		// The unsigned int argument is converted to unsigned decimal notation. The
		// precision, if any, gives the minimum number of digits that must appear; if
		// the converted value requires fewer digits, it is padded on the left with
		// zeros.  The default precision is 1.  When 0 is printed with an explicit
		// precision 0, the output is empty.
		format++
		var arg uint64
		if isWindows && mod == modL {
			mod = modNone
		}
		switch mod {
		case modNone:
			arg = uint64(VaUint32(args))
		case modL, modLL, mod64:
			arg = VaUint64(args)
		case modH:
			arg = uint64(uint16(VaInt32(args)))
		case modHH:
			arg = uint64(uint8(VaInt32(args)))
		case mod32:
			arg = uint64(VaInt32(args))
		case modZ:
			arg = uint64(VaInt64(args))
		default:
			panic(todo("", mod))
		}

		str = printfInteger(false, arg, c, flags, width, intPrec(hasPrecision, prec))
	case 'o':
		// The unsigned int argument is converted to unsigned octal notation. The
		// precision, if any, gives the minimum number of digits that must appear; if
		// the converted value requires fewer digits, it is padded on the left with
		// zeros.  The default precision is 1.  When 0 is printed with an explicit
		// precision 0, the output is empty.
		format++
		var arg uint64
		if isWindows && mod == modL {
			mod = modNone
		}
		switch mod {
		case modNone:
			arg = uint64(VaUint32(args))
		case modL, modLL, mod64:
			arg = VaUint64(args)
		case modH:
			arg = uint64(uint16(VaInt32(args)))
		case modHH:
			arg = uint64(uint8(VaInt32(args)))
		case mod32:
			arg = uint64(VaInt32(args))
		default:
			panic(todo("", mod))
		}

		str = printfInteger(false, arg, c, flags, width, intPrec(hasPrecision, prec))
	case 'b':
		// Base 2.
		format++
		var arg uint64
		if isWindows && mod == modL {
			mod = modNone
		}
		switch mod {
		case modNone:
			arg = uint64(VaUint32(args))
		case modL, modLL, mod64:
			arg = VaUint64(args)
		case modH:
			arg = uint64(uint16(VaInt32(args)))
		case modHH:
			arg = uint64(uint8(VaInt32(args)))
		case mod32:
			arg = uint64(VaInt32(args))
		default:
			panic(todo("", mod))
		}

		if arg == 0 && hasPrecision && prec == 0 {
			break
		}

		if hasPrecision {
			panic(todo("", prec))
		}

		f := spec + "b"
		str = fmt.Sprintf(f, arg)
	case 'I':
		if !isWindows {
			panic(todo("%#U", c))
		}

		format++
		switch c = *(*byte)(unsafe.Pointer(format)); c {
		case 'x', 'X':
			// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-wsprintfa
			//
			// Ix, IX
			//
			// 64-bit unsigned hexadecimal integer in lowercase or uppercase on 64-bit
			// platforms, 32-bit unsigned hexadecimal integer in lowercase or uppercase on
			// 32-bit platforms.
			if unsafe.Sizeof(int(0)) == 4 {
				mod = mod32
			}
		case '3':
			// https://en.wikipedia.org/wiki/Printf_format_string#Length_field
			//
			// I32	For integer types, causes printf to expect a 32-bit (double word) integer argument.
			format++
			switch c = *(*byte)(unsafe.Pointer(format)); c {
			case '2':
				format++
				mod = mod32
				goto more
			default:
				panic(todo("%#U", c))
			}
		case '6':
			// https://en.wikipedia.org/wiki/Printf_format_string#Length_field
			//
			// I64	For integer types, causes printf to expect a 64-bit (quad word) integer argument.
			format++
			switch c = *(*byte)(unsafe.Pointer(format)); c {
			case '4':
				format++
				mod = mod64
				goto more
			default:
				panic(todo("%#U", c))
			}
		default:
			panic(todo("%#U", c))
		}
		fallthrough
	case 'X':
		fallthrough
	case 'x':
		// The unsigned int argument is converted to unsigned hexadecimal notation.
		// The letters abcdef are used for x  conversions;  the letters ABCDEF are used
		// for X conversions.  The precision, if any, gives the minimum number of
		// digits that must appear; if the converted value requires fewer digits, it is
		// padded on the left with zeros.  The default precision is 1.  When 0 is
		// printed with an explicit precision 0, the output is empty.
		format++
		var arg uint64
		if isWindows && mod == modL {
			mod = modNone
		}
		switch mod {
		case modNone:
			arg = uint64(VaUint32(args))
		case modL, modLL, mod64:
			arg = VaUint64(args)
		case modH:
			arg = uint64(uint16(VaInt32(args)))
		case modHH:
			arg = uint64(uint8(VaInt32(args)))
		case mod32:
			arg = uint64(VaInt32(args))
		case modZ:
			arg = uint64(VaInt64(args))
		default:
			panic(todo("", mod))
		}

		str = printfInteger(false, arg, c, flags, width, intPrec(hasPrecision, prec))
	case 'e', 'E':
		// The double argument is rounded and converted in the style [-]d.ddde±dd where
		// there is one digit before the decimal-point character and the number of
		// digits after it is equal to  the precision;  if the precision is missing, it
		// is taken as 6; if the precision is zero, no decimal-point character appears.
		// An E conversion uses the letter E (rather than e) to intro‐ duce the
		// exponent.  The exponent always contains at least two digits; if the value is
		// zero, the exponent is 00.
		format++
		arg := VaFloat64(args)
		if math.IsInf(arg, 0) || math.IsNaN(arg) {
			str = printfNonFinite(arg, flags, width, c == 'E')
			break
		}
		if !hasPrecision || prec < 0 {
			prec = 6
		}
		f := fmt.Sprintf("%s.%d%c", spec, prec, c)
		str = fmt.Sprintf(f, arg)
	case 'f', 'F':
		// The double argument is rounded and converted to decimal notation in the
		// style [-]ddd.ddd, where the number of digits after the decimal-point
		// character  is  equal  to  the  precision specification.   If  the  precision
		// is missing, it is taken as 6; if the precision is explicitly zero, no
		// decimal-point character appears.  If a decimal point appears, at least one
		// digit appears before it.
		format++
		arg := VaFloat64(args)
		if math.IsInf(arg, 0) || math.IsNaN(arg) {
			str = printfNonFinite(arg, flags, width, c == 'F')
			break
		}
		if !hasPrecision || prec < 0 {
			prec = 6
		}
		f := fmt.Sprintf("%s.%d%c", spec, prec, c)
		str = fmt.Sprintf(f, arg)
	case 'G':
		fallthrough
	case 'g':
		// The double argument is converted in style f or e (or F or E for G
		// conversions).  The precision specifies the number of significant digits.  If
		// the precision is missing, 6 digits are given;  if the precision is zero, it
		// is treated as 1.  Style e is used if the exponent from its conversion is
		// less than -4 or greater than or equal to the precision.  Trailing zeros are
		// removed from the fractional part of the result; a decimal point appears only
		// if it is followed by at least one digit.
		format++
		arg := VaFloat64(args)
		if math.IsInf(arg, 0) || math.IsNaN(arg) {
			str = printfNonFinite(arg, flags, width, c == 'G')
			break
		}
		if !hasPrecision || prec < 0 {
			prec = 6
		}
		if prec == 0 {
			prec = 1
		}

		f := fmt.Sprintf("%s.%d%c", spec, prec, c)
		str = fmt.Sprintf(f, arg)
	case 'a', 'A':
		// For a conversion, the double argument is converted to hexadecimal
		// notation (using the letters abcdef) in the style [-]0xh.hhhhp±d; for A
		// conversion the prefix 0X, the letters ABCDEF, and the exponent
		// separator P is used.  There is one hexadecimal digit before the
		// decimal point, and the number of digits after it is equal to the
		// precision.  The default precision suffices for an exact
		// representation of the value if an exact representation in base 2
		// exists and otherwise is sufficiently large to distinguish values of
		// type double.  The digit before the decimal point is unspecified for
		// nonnormalized numbers, and nonzero but otherwise unspecified for
		// normalized numbers.  The exponent always contains at least one
		// digit; if the value is zero, the exponent is 0.
		format++
		if !hasWidth {
			width = 0
		}
		if !hasPrecision {
			prec = -1
		}
		str = printfHexFloat(VaFloat64(args), flags, width, prec, c == 'A')
	case 's':
		// If  no l modifier is present: the const char * argument is expected to be a
		// pointer to an array of character type (pointer to a string).  Characters
		// from the array are written up to (but not including) a terminating null byte
		// ('\0'); if a precision is specified, no more than the number specified are
		// written.  If a precision  is  given,  no  null  byte  need  be present; if
		// the precision is not specified, or is greater than the size of the array,
		// the array must contain a terminating null byte.
		//
		// If  an  l  modifier  is  present: the const wchar_t * argument is expected
		// to be a pointer to an array of wide characters.  Wide characters from the
		// array are converted to multibyte characters (each by a call to the
		// wcrtomb(3) function, with a conversion state starting in the initial state
		// before the first wide character), up to and including a terminating null
		// wide  character.   The  resulting  multibyte  characters are written up to
		// (but not including) the terminating null byte.  If a precision is specified,
		// no more bytes than the number specified are written, but no partial
		// multibyte characters are written.  Note that the precision determines the
		// number of bytes written, not the number of wide characters or  screen
		// positions.   The  array  must contain a terminating null wide character,
		// unless a precision is given and it is so small that the number of bytes
		// written exceeds it before the end of the array is reached.
		format++
		arg := VaUintptr(args)
		switch mod {
		case modNone:
			// Like musl: bytes, not runes, a NULL pointer is "(null)", and
			// the '0' flag doesn't apply.
			var s string
			switch prec := intPrec(hasPrecision, prec); {
			case arg == 0:
				s = "(null)"
				if prec >= 0 && prec < len(s) {
					s = s[:prec]
				}
			case prec >= 0:
				n := 0
				for n < prec && *(*byte)(unsafe.Pointer(arg + uintptr(n))) != 0 {
					n++
				}
				s = string(unsafe.Slice((*byte)(unsafe.Pointer(arg)), n))
			default:
				s = GoString(arg)
			}
			str = printfPad("", "", s, flags, width, false)
		default:
			panic(todo(""))
		}
	case 'p':
		// The void * pointer argument is printed in hexadecimal (as if by %#x or
		// %#lx).
		format++
		var s string
		switch runtime.GOOS {
		case "windows":
			switch runtime.GOARCH {
			case "386", "arm":
				s = fmt.Sprintf("%08X", VaUintptr(args))
			default:
				s = fmt.Sprintf("%016X", VaUintptr(args))
			}
		case "illumos":
			// Like the system libc: no 0x prefix, NULL is 0.
			s = fmt.Sprintf("%x", VaUintptr(args))
		default:
			s = fmt.Sprintf("%#0x", VaUintptr(args))
		}
		// Of the flags, only '-' applies: C leaves the others undefined
		// for %p, and the C libraries pad it with spaces. A negative
		// width from '*' put a '-' into spec too.
		if width < 0 {
			width = -width
		}
		switch {
		case !hasWidth:
			str = s
		case strings.Contains(spec, "-"):
			str = fmt.Sprintf("%-*s", width, s)
		default:
			str = fmt.Sprintf("%*s", width, s)
		}
	case 'c':
		// If no l modifier is present, the int argument is converted to an unsigned
		// char, and the resulting character is written.  If an l modifier is present,
		// the wint_t (wide character) ar‐ gument is converted to a multibyte sequence
		// by a call to the wcrtomb(3) function, with a conversion state starting in
		// the initial state, and the resulting multibyte string is  writ‐ ten.
		format++
		switch mod {
		case modNone:
			// The width applies, the '0' flag doesn't, as in musl.
			str = printfPad("", "", string([]byte{byte(VaInt32(args))}), flags, width, false)
		default:
			panic(todo(""))
		}
	case '%':
		// A '%' is written.  No argument is converted.  The complete conversion
		// specification is '%%'.
		format++
		buf.WriteByte('%')
	default:
		panic(todo("%#U", c))
	}

	buf.WriteString(str)
	return format
}

// Field width
//
// An optional decimal digit string (with nonzero first digit) specifying a
// minimum field width.  If the converted value has fewer characters than the
// field width, it will be padded with spa‐ ces on the left (or right, if the
// left-adjustment flag has been given).  Instead of a decimal digit string one
// may write "*" or "*m$" (for some decimal integer m) to specify that the
// field width  is  given  in the next argument, or in the m-th argument,
// respectively, which must be of type int.  A negative field width is taken as
// a '-' flag followed by a positive field width.  In no case does a
// nonexistent or small field width cause truncation of a field; if the result
// of a conversion is wider than the field width, the field is expanded to
// contain the conversion result.
func parseFieldWidth(format uintptr, args *uintptr) (_ uintptr, n int, ok bool) {
	first := true
	for {
		var digit int
		switch c := *(*byte)(unsafe.Pointer(format)); {
		case first && c == '0':
			return format, n, ok
		case first && c == '*':
			format++
			switch c := *(*byte)(unsafe.Pointer(format)); {
			case c >= '0' && c <= '9':
				panic(todo(""))
			default:
				return format, int(VaInt32(args)), true
			}
		case c >= '0' && c <= '9':
			format++
			ok = true
			first = false
			digit = int(c) - '0'
		default:
			return format, n, ok
		}

		n0 := n
		n = 10*n + digit
		if n < n0 {
			panic(todo(""))
		}
	}
}

// Precision
//
// An  optional precision, in the form of a period ('.')  followed by an
// optional decimal digit string.  Instead of a decimal digit string one may
// write "*" or "*m$" (for some decimal integer m) to specify that the
// precision is given in the next argument, or in the m-th argument,
// respectively, which must be of type int.  If the precision is given as just
// '.', the  precision  is taken  to  be  zero.  A negative precision is taken
// as if the precision were omitted.  This gives the minimum number of digits
// to appear for d, i, o, u, x, and X conversions, the number of digits to
// appear after the radix character for a, A, e, E, f, and F conversions, the
// maximum number of significant digits for g and G conversions, or the maximum
// number of characters to be printed from a string for s and S conversions.
func parsePrecision(format uintptr, args *uintptr) (_ uintptr, n int, ok bool) {
	for {
		switch c := *(*byte)(unsafe.Pointer(format)); c {
		case '.':
			format++
			first := true
			for {
				switch c := *(*byte)(unsafe.Pointer(format)); {
				case first && c == '*':
					format++
					n = int(VaInt32(args))
					return format, n, true
				case c >= '0' && c <= '9':
					format++
					first = false
					n0 := n
					n = 10*n + (int(c) - '0')
					if n < n0 {
						panic(todo(""))
					}
				default:
					return format, n, true
				}
			}
		default:
			return format, 0, false
		}
	}
}

// Length modifier
//
// Here, "integer conversion" stands for d, i, o, u, x, or X conversion.
//
// hh     A following integer conversion corresponds to a signed char or
// unsigned char argument, or a following n conversion corresponds to a pointer
// to a signed char argument.
//
// h      A following integer conversion corresponds to a short int or unsigned
// short int argument, or a following n conversion corresponds to a pointer to
// a short int argument.
//
// l      (ell)  A following integer conversion corresponds to a long int or
// unsigned long int argument, or a following n conversion corresponds to a
// pointer to a long int argument, or a fol‐ lowing c conversion corresponds to
// a wint_t argument, or a following s conversion corresponds to a pointer to
// wchar_t argument.
//
// ll     (ell-ell).  A following integer conversion corresponds to a long long
// int or unsigned long long int argument, or a following n conversion
// corresponds to a pointer to a long long int argument.
//
// q      A synonym for ll.  This is a nonstandard extension, derived from BSD;
// avoid its use in new code.
//
// L      A following a, A, e, E, f, F, g, or G conversion corresponds to a
// long double argument.  (C99 allows %LF, but SUSv2 does not.)
//
// j      A following integer conversion corresponds to an intmax_t or
// uintmax_t argument, or a following n conversion corresponds to a pointer to
// an intmax_t argument.
//
// z      A following integer conversion corresponds to a size_t or ssize_t
// argument, or a following n conversion corresponds to a pointer to a size_t
// argument.
//
// Z      A nonstandard synonym for z that predates the appearance of z.  Do
// not use in new code.
//
// t      A following integer conversion corresponds to a ptrdiff_t argument,
// or a following n conversion corresponds to a pointer to a ptrdiff_t
// argument.

func parseLengthModifier(format uintptr) (_ uintptr, n int) {
	switch c := *(*byte)(unsafe.Pointer(format)); c {
	case 'h':
		format++
		n = modH
		switch c := *(*byte)(unsafe.Pointer(format)); c {
		case 'h':
			format++
			n = modHH
		}
		return format, n
	case 'l':
		format++
		n = modL
		switch c := *(*byte)(unsafe.Pointer(format)); c {
		case 'l':
			format++
			n = modLL
		}
		return format, n
	case 'q':
		panic(todo(""))
	case 'L':
		format++
		n = modLD
		return format, n
	case 'j':
		format++
		n = modJ
		return format, n
	case 'z':
		format++
		return format, modZ
	case 'Z':
		format++
		return format, modCapitalZ
	case 't':
		format++
		return format, modT
	default:
		return format, 0
	}
}

// intPrec returns the precision of a conversion, or -1 if there is none. A
// negative precision from '*' counts as none.
func intPrec(hasPrecision bool, prec int) int {
	if !hasPrecision || prec < 0 {
		return -1
	}

	return prec
}

// printfInteger formats %d, %i, %u, %o, %x or %X like musl does. u is the
// magnitude, and neg is set for a negative %d or %i. A precision (prec >= 0)
// gives the minimum number of digits, and the '0' flag then doesn't apply.
// Precision 0 prints no digits for 0, except for %#o, which always starts
// with a 0. '#' prefixes a nonzero %x or %X with 0x or 0X. '+' and ' ' apply
// only to %d and %i. A negative width adjusts left, as the '-' flag does.
func printfInteger(neg bool, u uint64, conv byte, flags string, width, prec int) string {
	var digits string
	switch {
	case u == 0:
		// No digits: the precision decides, see below.
	case conv == 'o':
		digits = strconv.FormatUint(u, 8)
	case conv == 'x':
		digits = strconv.FormatUint(u, 16)
	case conv == 'X':
		digits = strings.ToUpper(strconv.FormatUint(u, 16))
	default:
		digits = strconv.FormatUint(u, 10)
	}
	var sign string
	switch conv {
	case 'd', 'i':
		switch {
		case neg:
			sign = "-"
		case strings.Contains(flags, "+"):
			sign = "+"
		case strings.Contains(flags, " "):
			sign = " "
		}
	case 'x', 'X':
		if u != 0 && strings.Contains(flags, "#") {
			sign = "0" + string(conv)
		}
	}
	zeroOK := prec < 0
	if conv == 'o' && strings.Contains(flags, "#") && prec < len(digits)+1 {
		prec = len(digits) + 1
	}
	switch {
	case u == 0 && prec == 0:
		// Nothing.
	case u == 0:
		prec = max(prec, 1)
	default:
		prec = max(prec, len(digits))
	}
	if prec > len(digits) {
		digits = strings.Repeat("0", prec-len(digits)) + digits
	}
	return printfPad(sign, "", digits, flags, width, zeroOK)
}

// printfNonFinite formats an infinity or a NaN for %a, %e, %f and %g like
// musl does: inf or nan, or INF or NAN if upper is set, with the sign of x or
// the sign the '+' or ' ' flag asks for, and padded with spaces even with the
// '0' flag. A negative width adjusts left, as the '-' flag does.
func printfNonFinite(x float64, flags string, width int, upper bool) string {
	s := "inf"
	if math.IsNaN(x) {
		s = "nan"
	}
	if upper {
		s = strings.ToUpper(s)
	}
	return printfPad(printfSign(x, flags), "", s, flags, width, false)
}

// printfSign returns the sign that %a, %e, %f and %g write before x.
func printfSign(x float64, flags string) string {
	switch {
	case math.Signbit(x):
		return "-"
	case strings.Contains(flags, "+"):
		return "+"
	case strings.Contains(flags, " "):
		return " "
	default:
		return ""
	}
}

// printfPad pads sign+prefix+body to width: on the right with spaces if the
// '-' flag is set or width is negative, otherwise on the left with zeros
// after the prefix if zeroOK and the '0' flag are set, otherwise on the left
// with spaces.
func printfPad(sign, prefix, body, flags string, width int, zeroOK bool) string {
	left := strings.Contains(flags, "-")
	if width < 0 {
		left = true
		width = -width
	}
	n := width - len(sign) - len(prefix) - len(body)
	switch {
	case n <= 0:
		return sign + prefix + body
	case left:
		return sign + prefix + body + strings.Repeat(" ", n)
	case zeroOK && strings.Contains(flags, "0"):
		return sign + prefix + strings.Repeat("0", n) + body
	default:
		return strings.Repeat(" ", n) + sign + prefix + body
	}
}

// printfHexFloat formats x for %a, or for %A if upper is set, like musl does.
// A nonzero x has the leading digit 1, subnormals included. Without a
// precision (prec < 0) the fraction has as many digits as x needs, otherwise
// x is rounded half to even to prec digits, and the leading digit can then
// become 2. The exponent has only the digits it needs. A negative width
// adjusts left, as the '-' flag does.
func printfHexFloat(x float64, flags string, width, prec int, upper bool) string {
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return printfNonFinite(x, flags, width, upper)
	}

	digits := "0123456789abcdef"
	prefix := "0x"
	expSep := "p"
	if upper {
		digits = "0123456789ABCDEF"
		prefix = "0X"
		expSep = "P"
	}
	b := math.Float64bits(x)
	e := int(b>>52) & 0x7ff
	m := b & (1<<52 - 1)
	switch {
	case e == 0 && m == 0:
		// Zero has the exponent 0.
	case e == 0:
		n := bits.LeadingZeros64(m) - 11
		m <<= n
		e = -1022 - n
	default:
		m |= 1 << 52
		e -= 1023
	}
	// m is now the leading digit followed by 13 hexadecimal digits of
	// fraction.
	nd := 13
	if prec >= 0 && prec < nd {
		drop := uint(4 * (nd - prec))
		rem := m & (1<<drop - 1)
		m >>= drop
		if half := uint64(1) << (drop - 1); rem > half || rem == half && m&1 != 0 {
			m++
		}
		nd = prec
	}
	frac := make([]byte, nd, max(nd, prec))
	for i := nd - 1; i >= 0; i-- {
		frac[i] = digits[m&15]
		m >>= 4
	}
	switch {
	case prec < 0:
		frac = bytes.TrimRight(frac, "0")
	default:
		for len(frac) < prec {
			frac = append(frac, '0')
		}
	}
	var sb strings.Builder
	sb.WriteByte(digits[m])
	if len(frac) != 0 || strings.Contains(flags, "#") {
		sb.WriteByte('.')
	}
	sb.Write(frac)
	sb.WriteString(expSep)
	switch {
	case e < 0:
		sb.WriteByte('-')
		e = -e
	default:
		sb.WriteByte('+')
	}
	sb.WriteString(strconv.Itoa(e))
	return printfPad(printfSign(x, flags), prefix, sb.String(), flags, width, true)
}
