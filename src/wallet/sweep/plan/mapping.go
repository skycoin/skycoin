// Package plan src/wallet/sweep/plan/mapping.go c1-coin-wallet
package plan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseSource reads a source written as an account number or "account:N", else
// an address, optionally written "address:X".
func ParseSource(s string) (Source, error) {
	s = strings.TrimSpace(s)
	if v, ok := strings.CutPrefix(s, "address:"); ok {
		return Source{Address: v}, nil
	}
	if v, ok := strings.CutPrefix(s, "account:"); ok {
		s = v
	}
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		a := uint32(n)
		return Source{Account: &a}, nil
	}
	if s == "" {
		return Source{}, errors.New("the source is empty")
	}
	return Source{Address: s}, nil
}

// ParseDestination reads a destination written as "account:N", "next-account",
// "wallet" for the wallet itself, "wallet:ID", or an address, optionally
// written "address:X".
func ParseDestination(s string) (Destination, error) {
	s = strings.TrimSpace(s)
	switch {
	case s == "next-account":
		return Destination{Kind: DestAccount}, nil
	case s == "wallet":
		return Destination{Kind: DestWallet}, nil
	case strings.HasPrefix(s, "account:"):
		n, err := strconv.ParseUint(strings.TrimPrefix(s, "account:"), 10, 32)
		if err != nil {
			return Destination{}, fmt.Errorf("bad account in %q", s)
		}
		a := uint32(n)
		return Destination{Kind: DestAccount, Account: &a}, nil
	case strings.HasPrefix(s, "wallet:"):
		return Destination{Kind: DestWallet, Wallet: strings.TrimPrefix(s, "wallet:")}, nil
	case s == "":
		return Destination{}, errors.New("the destination is empty")
	}
	return Destination{Kind: DestAddress, Address: strings.TrimPrefix(s, "address:")}, nil
}

// ParseMapping reads "<source>=<destination>".
func ParseMapping(s string) (Mapping, error) {
	from, to, ok := strings.Cut(s, "=")
	if !ok {
		return Mapping{}, fmt.Errorf("%q: want <source>=<destination>", s)
	}
	src, err := ParseSource(from)
	if err != nil {
		return Mapping{}, fmt.Errorf("%q: %w", s, err)
	}
	dst, err := ParseDestination(to)
	if err != nil {
		return Mapping{}, fmt.Errorf("%q: %w", s, err)
	}
	return Mapping{From: src, To: dst}, nil
}

// String writes s as ParseSource reads it.
func (s Source) String() string {
	if s.Account != nil {
		return strconv.FormatUint(uint64(*s.Account), 10)
	}
	return s.Address
}

// String writes d as ParseDestination reads it.
func (d Destination) String() string {
	switch d.Kind {
	case DestAccount:
		if d.Account == nil {
			return "next-account"
		}
		return "account:" + strconv.FormatUint(uint64(*d.Account), 10)
	case DestWallet:
		if d.Wallet == "" {
			return "wallet"
		}
		return "wallet:" + d.Wallet
	}
	return d.Address
}

// String writes m as ParseMapping reads it.
func (m Mapping) String() string {
	return m.From.String() + "=" + m.To.String()
}
