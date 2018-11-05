// Package useragent implements methods for managing Skycoin user agents.
//
// A skycoin user agent has the following format:
//
//   `$NAME:$VERSION[$GIT_HASH]($REMARK)`
//
// `$NAME` and `$VERSION` are required.
//
// * `$NAME` is the coin or application's name, e.g. `Skycoin`. It can contain the following characters: `A-Za-z0-9\-_+`.
// * `$VERSION` must be a valid [semver](http://semver.org/) version, e.g. `1.2.3` or `1.2.3-rc1`.
//   Semver has the option of including build metadata such as the git commit hash, but this is not included by the default client.
// * `$REMARK` is optional. If not present, the enclosing brackets `()` should be omitted.
//   It can contain the following characters: `A-Za-z0-9\-_+;:!$%,.=?~ ` (includes the space character).
package useragent

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

const (
	// IllegalChars are printable ascii characters forbidden from a user agent string. All other ascii or bytes are also forbidden.
	IllegalChars = `<>&"'#@|{}` + "`"
	// MaxLen the maximum length of a user agent
	MaxLen = 256

	// NamePattern is the regex pattern for the name portion of the user agent
	NamePattern = `[A-Za-z0-9\-_+]+`
	// VersionPattern is the regex pattern for the version portion of the user agent
	VersionPattern = `[0-9]+\.[0-9]+\.[0-9][A-Za-z0-9\-.+]*`
	// RemarkPattern is the regex pattern for the remark portion of the user agent
	RemarkPattern = `[A-Za-z0-9\-_+;:!$%,.=?~ ]+`

	// Pattern is the regex pattern for the user agent in entirety
	Pattern = `^(` + NamePattern + `):(` + VersionPattern + `)(\(` + RemarkPattern + `\))?$`
)

var (
	illegalCharsSanitizeRe *regexp.Regexp
	illegalCharsCheckRe    *regexp.Regexp
	re                     *regexp.Regexp

	// ErrIllegalChars user agent contains illegal characters
	ErrIllegalChars = errors.New("User agent has invalid character(s)")
	// ErrTooLong user agent exceeds a certain max length
	ErrTooLong = errors.New("User agent is too long")
	// ErrMalformed user agent does not match the user agent pattern
	ErrMalformed = errors.New("User agent is malformed")
	// ErrEmpty user agent is an empty string
	ErrEmpty = errors.New("User agent is an empty string")
)

func init() {
	illegalCharsSanitizeRe = regexp.MustCompile(fmt.Sprintf("([^[:print:]]|[%s])+", IllegalChars))
	illegalCharsCheckRe = regexp.MustCompile(fmt.Sprintf("[^[:print:]]|[%s]", IllegalChars))
	re = regexp.MustCompile(Pattern)
}

// Data holds parsed user agent data
type Data struct {
	Coin    string
	Version string
	Remark  string
}

// Empty returns true if the Data is empty
func (d Data) Empty() bool {
	return d == (Data{})
}

// Build builds a user agent string. Returns an error if the user agent would be invalid.
func (d Data) Build() (string, error) {
	if d.Coin == "" {
		return "", errors.New("missing coin name")
	}
	if d.Version == "" {
		return "", errors.New("missing version")
	}

	_, err := semver.Parse(d.Version)
	if err != nil {
		return "", err
	}

	s := d.build()

	if err := validate(s); err != nil {
		return "", err
	}

	d2, err := Parse(s)
	if err != nil {
		return "", fmt.Errorf("Built a user agent that fails to parse: %q %v", s, err)
	}

	if d2 != d {
		return "", errors.New("Built a user agent that does not parse to the original format")
	}

	return s, nil
}

// MustBuild calls Build and panics on error
func (d Data) MustBuild() string {
	s, err := d.Build()
	if err != nil {
		panic(err)
	}
	return s
}

func (d Data) build() string {
	if d.Coin == "" || d.Version == "" {
		return ""
	}

	remark := d.Remark
	if remark != "" {
		remark = fmt.Sprintf("(%s)", remark)
	}

	return fmt.Sprintf("%s:%s%s", d.Coin, d.Version, remark)
}

// MarshalJSON marshals Data as JSON
func (d Data) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.build())), nil
}

// UnmarshalJSON unmarshals []byte to Data
func (d *Data) UnmarshalJSON(v []byte) error {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// Parse parses a user agent string to Data
func Parse(userAgent string) (Data, error) {
	if len(userAgent) == 0 {
		return Data{}, ErrEmpty
	}

	if err := validate(userAgent); err != nil {
		return Data{}, err
	}

	subs := re.FindAllStringSubmatch(userAgent, -1)

	if len(subs) == 0 {
		return Data{}, ErrMalformed
	}

	m := subs[0]

	if m[0] != userAgent {
		// This should not occur since the pattern has ^$ boundaries applied, but just in case
		return Data{}, errors.New("User agent did not match pattern completely")
	}

	coin := m[1]
	version := m[2]
	remark := m[3]

	if _, err := semver.Parse(version); err != nil {
		return Data{}, fmt.Errorf("User agent version is not a valid semver: %v", err)
	}

	remark = strings.TrimPrefix(remark, "(")
	remark = strings.TrimSuffix(remark, ")")

	return Data{
		Coin:    coin,
		Version: version,
		Remark:  remark,
	}, nil
}

// MustParse parses and panics on error
func MustParse(userAgent string) Data {
	d, err := Parse(userAgent)
	if err != nil {
		panic(err)
	}

	return d
}

// validate validates a user agent string. The user agent must not contain illegal characters.
func validate(userAgent string) error {
	if len(userAgent) > MaxLen {
		return ErrTooLong
	}

	if illegalCharsCheckRe.MatchString(userAgent) {
		return ErrIllegalChars
	}

	return nil
}

// Sanitize removes illegal characters from a user agent string
func Sanitize(userAgent string) string {
	return illegalCharsSanitizeRe.ReplaceAllLiteralString(userAgent, "")
}
