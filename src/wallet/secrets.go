package wallet

import "encoding/json"

// secrets key name
const (
	secretSeed     = "seed"
	secretLastSeed = "lastSeed"
)

type secrets map[string]string

func (s secrets) get(key string) (string, bool) {
	v, ok := s[key]
	return v, ok
}

func (s secrets) set(key, v string) {
	s[key] = v
}

func (s secrets) serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s secrets) deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s secrets) erase() {
	for k := range s {
		s[k] = ""
		delete(s, k)
	}
}
