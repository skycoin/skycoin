// Package ui2 src/ui2/reset.go c3-ui-wallet
package ui2

import (
	"errors"

	"github.com/skycoin/skycoin/src/api"
)

var errResetSeed = errors.New("create.err-seed")

// ResetPassword replaces the password of encrypted wallet id using its
// seed, or removes it when the new password is empty, as the Angular
// wallet's reset page does. It returns the notice for the wallet page.
func ResetPassword(n Node, id string, post func(string) string) (string, error) {
	seed := normalizeSeed(post("seed"))
	pw := post("password")
	switch {
	case seed == "":
		return "", errResetSeed
	case pw != post("password2"):
		return "", errPasswordTwo
	}
	if _, err := n.RecoverWallet(api.WalletRecoverRequest{ID: id, Seed: seed, SeedPassphrase: post("passphrase"), Password: pw}); err != nil {
		return "", err
	}
	if pw == "" {
		return "decrypted", nil
	}
	return "reset", nil
}
