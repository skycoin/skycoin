package wallet

import "encoding/json"

const (
	secretSeed           = "seed"
	secretLastSeed       = "lastSeed"
	secretSeedPassphrase = "seedPassphrase"
)

// Secrets hold secret data, to be encrypted
type Secrets map[string]string

func (s Secrets) get(key string) (string, bool) {
	v, ok := s[key]
	return v, ok
}

func (s Secrets) set(key, v string) {
	s[key] = v
}

func (s Secrets) serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s Secrets) deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s Secrets) erase() {
	for k := range s {
		s[k] = ""
		delete(s, k)
	}
}
