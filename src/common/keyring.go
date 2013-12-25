package sb

import (
    "os"
    "fmt"
    "errors"
    "path/filepath"
)

type Keypair struct {
    Name string
    Public []byte
    Private []byte
    }

    func NewKeypair(name string) *Keypair {
        private, public := GenerateKeyPair()
        return &Keypair{Name: name, Public: public, Private: private}
    }


type Keyring struct {
    Location string
    Keys map[string]*Keypair
    }

    func NewKeyring(location string) *Keyring {
        path, err := setupKeyringLocation(location + ".keyring")
        if err != nil {
            fmt.Printf("Failed to setup keyring location: %s\n", err.Error())
            return nil
        }
        keyring := &Keyring{Location: path, Keys: make(map[string]*Keypair, 4)}
        keyring.Load()
        return keyring
    }

    func (self *Keyring) Load() error {
        err := LoadGob(self.Location, self)
        if err != nil {
            fmt.Printf("Failed to load keyring from %s. Reason: %v\n", self.Location, err)
            return err
        }
        fmt.Printf("Loaded keyring from %s\n", self.Location)
        return nil
    }

    func (self *Keyring) Save() error {
        err := SaveGob(self.Location, self)
        if err != nil {
            fmt.Printf("Failed to save keyring to %s. Reason: %v\n", self.Location, err)
            return err
        }
        fmt.Printf("Saved keyring to %s\n", self.Location)
        return nil
    }

    func (self *Keyring) Create(name string) *Keypair {
        kp := NewKeypair(name)
        self.Keys[name] = kp
        self.Save()
        fmt.Printf("Created and saved new keypair \"%s\" at %s\n", name, self.Location)
        return kp
    }

    func (self *Keyring) Remove(name string) {
        delete(self.Keys, name)
    }

    func (self *Keyring) Get(name string) *Keypair {
        return self.Keys[name]
    }

    func (self *Keyring) Empty() bool {
        return (len(self.Keys) == 0)
    }


func setupKeyringLocation(location string) (string, error) {
    if location == "" {
        return "", errors.New("Keyring location cannot be empty")
    }
    path := filepath.Join(DataDirectory, "keys/")
    perms := os.FileMode(0755)
    err := os.MkdirAll(path, perms)
    path = filepath.Join(path, location)
    if err != nil {
        return "", err
    }
    return path, nil
}
