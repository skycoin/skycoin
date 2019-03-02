package note

import (
	"fmt"
	"os"
	"sync"
)

// Config is note service config.
type Config struct {
	NotesDir      string
	EnableNoteAPI bool
}

// NewConfig creates a default config.
func NewConfig() Config {
	return Config{
		NotesDir: "./notes/",
	}
}

// Manager notes manager struct
type Manager struct {
	sync.RWMutex
	config Config
	notes  map[string]string
}

// NewManager creates a new note manager
func NewManager(c Config) (*Manager, error) {
	if err := os.MkdirAll(c.NotesDir, os.FileMode(0700)); err != nil {
		return nil, fmt.Errorf("failed to create note directory %s: %v", c.NotesDir, err)
	}

	n, err := LoadNotes(c.NotesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load all notes: %v", err)
	}

	m := &Manager{
		config: c,
		notes:  n,
	}

	return m, nil
}

// GetNotes returns a map with all notes
func (m *Manager) GetNotes() (map[string]string, error) {
	m.RLock()
	defer m.RUnlock()

	if !m.config.EnableNoteAPI {
		return nil, ErrNoteAPIDisabled
	}

	// copy m.notes to avoid its modification outside the function
	notes := make(map[string]string, len(m.notes))
	for k, v := range m.notes {
		notes[k] = v
	}

	return notes, nil
}

// GetNote returns a note by transaction ID
func (m *Manager) GetNote(txID string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	if !m.config.EnableNoteAPI {
		return "", ErrNoteAPIDisabled
	}

	if err := validateTxID(txID); err != nil {
		return "", ErrInvalidTxID
	}

	note, ok := m.notes[txID]
	if !ok {
		return "", ErrNoteNotExist
	}

	return note, nil
}

// AddNote adds a note for a transaction ID
func (m *Manager) AddNote(txID, note string) error {
	m.Lock()
	defer m.Unlock()

	if !m.config.EnableNoteAPI {
		return ErrNoteAPIDisabled
	}

	if err := validateTxID(txID); err != nil {
		return ErrInvalidTxID
	}

	if err := Save(m.config.NotesDir, txID, note); err != nil {
		return err
	}

	m.notes[txID] = note

	return nil
}

// RemoveNote removes a note by transaction ID
func (m *Manager) RemoveNote(txID string) error {
	m.Lock()
	defer m.Unlock()

	if !m.config.EnableNoteAPI {
		return ErrNoteAPIDisabled
	}

	if err := validateTxID(txID); err != nil {
		return ErrInvalidTxID
	}

	if _, ok := m.notes[txID]; !ok {
		return ErrNoteNotExist
	}

	if err := Remove(m.config.NotesDir, txID); err != nil {
		return err
	}

	delete(m.notes, txID)

	return nil
}
