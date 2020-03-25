package wallet

import (
	"bytes"
	"encoding/json"
)

type bip44WalletJSONDecoder struct {
}

func (d bip44WalletJSONDecoder) Encode(w *Bip44WalletNew) ([]byte, error) {
	rw := NewReadableBip44WalletNew(w)
	return json.MarshalIndent(rw, "", "    ")
}

func (d bip44WalletJSONDecoder) Decode(b []byte) (*Bip44WalletNew, error) {
	br := bytes.NewReader(b)
	rw := ReadableBip44WalletNew{}
	if err := json.NewDecoder(br).Decode(&rw); err != nil {
		return nil, err
	}
	return rw.ToWallet()
}

// ReadableBip44WalletNew readable bip44 wallet
type ReadableBip44WalletNew struct {
	Meta     `json:"meta"`
	Accounts ReadableBip44Accounts `json:"accounts"`
}

// NewReadableBip44WalletNew creates a readable bip44 wallet
func NewReadableBip44WalletNew(w *Bip44WalletNew) *ReadableBip44WalletNew {
	return &ReadableBip44WalletNew{
		Meta:     w.Meta.clone(),
		Accounts: w.accounts.ToReadable(),
	}
}

// ToWallet converts the readable bip44 wallet to a bip44 wallet
func (rw ReadableBip44WalletNew) ToWallet() (*Bip44WalletNew, error) {
	w := Bip44WalletNew{
		Meta: rw.Meta.clone(),
	}
	return &w, nil
}
