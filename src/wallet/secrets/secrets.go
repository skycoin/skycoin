package secrets

import "encoding/json"

const (
	// SecretSeed key of seed in Secrets
	SecretSeed = "seed"
	// SecretLastSeed key of last sees in Secrets
	SecretLastSeed = "lastSeed"
	// SecretSeedPassphrase key of seed passphrase in Secrets
	SecretSeedPassphrase = "seedPassphrase"
)

// Secrets hold secret data, to be encrypted
type Secrets map[string]string

// Get returns the secret value of given key
func (s Secrets) Get(key string) (string, bool) {
	v, ok := s[key]
	return v, ok
}

// Set sets the secret key and value
func (s Secrets) Set(key, v string) {
	s[key] = v
}

// Serialize encodes the secrets into []bytes
func (s Secrets) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

// Deserialize decodes  the secrets from []bytes
func (s Secrets) Deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}

// Erase wipes all secrets
func (s Secrets) Erase() {
	for k := range s {
		s[k] = ""
		delete(s, k)
	}
}
