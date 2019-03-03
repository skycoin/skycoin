package note

import (
	"fmt"
	"os"
	"testing"

	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/file"

	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	type expect struct {
		notes map[string]string
		err   error
	}

	tt := []struct {
		name   string
		config Config
		expect expect
	}{
		{
			name: "notespath: ./testdata/notes",
			config: Config{
				NotesDir:      "./testdata/notes",
				EnableNoteAPI: true,
			},
			expect: expect{
				notes: make(map[string]string),
				err:   nil,
			},
		},
		{
			name: "notespath: ./testdata",
			config: Config{
				NotesDir: "./testdata",
			},
			expect: expect{
				notes: map[string]string{
					"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
					"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
					"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
				},
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewManager(tc.config)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			_, err = os.Stat(tc.config.NotesDir)
			require.NoError(t, err)

			require.Equal(t, tc.expect.notes, m.notes)
		})
	}
}

func makeManager(t *testing.T, config Config) *Manager {
	m, err := NewManager(config)
	require.NoError(t, err)

	return m
}

func TestGetNotes(t *testing.T) {
	type expect struct {
		notes map[string]string
		err   error
	}

	tt := []struct {
		name   string
		config Config
		expect expect
	}{
		{
			name: "ok, no notes",
			config: Config{
				NotesDir:      "./testdata/notes",
				EnableNoteAPI: true,
			},
			expect: expect{
				notes: make(map[string]string),
				err:   nil,
			},
		},
		{
			name: "!ok: api disabled, no notes",
			config: Config{
				NotesDir: "./testdata/notes",
			},
			expect: expect{
				notes: nil,
				err:   ErrNoteAPIDisabled,
			},
		},
		{
			name: "ok, 3 notes",
			config: Config{
				NotesDir:      "./testdata",
				EnableNoteAPI: true,
			},
			expect: expect{
				notes: map[string]string{
					"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
					"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
					"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
				},
				err: nil,
			},
		},
		{
			name: "!ok: api disabled, 3 notes",
			config: Config{
				NotesDir: "./testdata",
			},
			expect: expect{
				notes: nil,
				err:   ErrNoteAPIDisabled,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := makeManager(t, tc.config)

			notes, err := m.GetNotes()
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.notes, notes)
		})
	}
}

func TestGetNote(t *testing.T) {
	type expect struct {
		note string
		err  error
	}

	noNotesEnabledManager := makeManager(t, Config{
		NotesDir:      "./testdata/notes",
		EnableNoteAPI: true,
	})
	noNotesManager := makeManager(t, Config{
		NotesDir: "./testdata/notes",
	})
	enabledManager := makeManager(t, Config{
		NotesDir:      "./testdata",
		EnableNoteAPI: true,
	})
	manager := makeManager(t, Config{
		NotesDir: "./testdata",
	})

	tt := []struct {
		name    string
		manager *Manager
		txID    string
		expect  expect
	}{
		{
			name:    "no notes, !ok: note does not exist, note: no",
			manager: noNotesEnabledManager,
			txID:    "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			expect: expect{
				note: "",
				err:  ErrNoteNotExist,
			},
		},
		{
			name:    "no notes, !ok: invalid txid, note: no",
			manager: noNotesEnabledManager,
			txID:    "txid1",
			expect: expect{
				note: "",
				err:  ErrInvalidTxID,
			},
		},
		{
			name:    "no notes, !ok: api disabled, note: no",
			manager: noNotesManager,
			txID:    "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			expect: expect{
				note: "",
				err:  ErrNoteAPIDisabled,
			},
		},
		{
			name:    "notes exist, !ok: api disabled, note: no",
			manager: manager,
			txID:    "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			expect: expect{
				note: "",
				err:  ErrNoteAPIDisabled,
			},
		},
		{
			name:    "notes exist, ok, note: note1",
			manager: enabledManager,
			txID:    "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			expect: expect{
				note: "note1",
				err:  nil,
			},
		},
		{
			name:    "notes exist, ok, note: note3",
			manager: enabledManager,
			txID:    "fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4",
			expect: expect{
				note: "note3",
				err:  nil,
			},
		},
		{
			name:    "notes exist, !ok: invalid txid, note: no",
			manager: enabledManager,
			txID:    "txid1",
			expect: expect{
				note: "",
				err:  ErrInvalidTxID,
			},
		},
		{
			name:    "notes exist, !ok: note does not exist, note: no",
			manager: enabledManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			expect: expect{
				note: "",
				err:  ErrNoteNotExist,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			note, err := tc.manager.GetNote(tc.txID)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.note, note)
		})
	}
}

func TestAddNote(t *testing.T) {
	noNotesEnabledManager := makeManager(t, Config{
		NotesDir:      "./testdata/notes",
		EnableNoteAPI: true,
	})
	noNotesManager := makeManager(t, Config{
		NotesDir: "./testdata/notes",
	})
	enabledManager := makeManager(t, Config{
		NotesDir:      "./testdata",
		EnableNoteAPI: true,
	})
	manager := makeManager(t, Config{
		NotesDir: "./testdata",
	})

	tt := []struct {
		name    string
		manager *Manager
		txID    string
		note    string
		expect  error
	}{
		{
			name:    "no notes, !ok: api disabled",
			manager: noNotesManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			note:    "note4",
			expect:  ErrNoteAPIDisabled,
		},
		{
			name:    "notes exist, !ok : api disabled",
			manager: manager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			note:    "note4",
			expect:  ErrNoteAPIDisabled,
		},
		{
			name:    "no notes, !ok: invalid tx id",
			manager: noNotesEnabledManager,
			txID:    "txid1",
			note:    "note4",
			expect:  ErrInvalidTxID,
		},
		{
			name:    "notes exist, !ok: invalid tx id",
			manager: enabledManager,
			txID:    "txid1",
			note:    "note4",
			expect:  ErrInvalidTxID,
		},
		{
			name:    "no notes, ok",
			manager: noNotesEnabledManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			note:    "note4",
			expect:  nil,
		},
		{
			name:    "notes exist, ok",
			manager: enabledManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			note:    "note4",
			expect:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.AddNote(tc.txID, tc.note)
			require.Equal(t, tc.expect, err)
			if err != nil {
				return
			}

			testutil.RequireFileExists(t, fmt.Sprintf("%s/%s.txnote", tc.manager.config.NotesDir, tc.txID))
		})
	}

	err := file.RemoveFile("./testdata/notes/b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d.txnote")
	require.NoError(t, err)
	err = file.RemoveFile("./testdata/b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d.txnote")
	require.NoError(t, err)
}

func TestRemoveNote(t *testing.T) {
	noNotesEnabledManager := makeManager(t, Config{
		NotesDir:      "./testdata/notes",
		EnableNoteAPI: true,
	})
	noNotesManager := makeManager(t, Config{
		NotesDir: "./testdata/notes",
	})
	enabledManager := makeManager(t, Config{
		NotesDir:      "./testdata",
		EnableNoteAPI: true,
	})
	manager := makeManager(t, Config{
		NotesDir: "./testdata",
	})

	tt := []struct {
		name    string
		manager *Manager
		txID    string
		expect  error
	}{
		{
			name:    "no notes, !ok: api disabled",
			manager: noNotesManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			expect:  ErrNoteAPIDisabled,
		},
		{
			name:    "notes exist, !ok: api disabled",
			manager: manager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			expect:  ErrNoteAPIDisabled,
		},
		{
			name:    "no notes, !ok: invalid txid",
			manager: noNotesEnabledManager,
			txID:    "txid1",
			expect:  ErrInvalidTxID,
		},
		{
			name:    "notes exist, !ok: invalid txid",
			manager: enabledManager,
			txID:    "txid1",
			expect:  ErrInvalidTxID,
		},
		{
			name:    "no notes, !ok: no such note",
			manager: noNotesEnabledManager,
			txID:    "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			expect:  ErrNoteNotExist,
		},
		{
			name:    "notes exist, !ok: no such note",
			manager: enabledManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			expect:  ErrNoteNotExist,
		},
		{
			name:    "ok",
			manager: noNotesEnabledManager,
			txID:    "b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
			expect:  nil,
		},
	}

	err := noNotesEnabledManager.AddNote("b2483c816a7b18a628b796def151aad61d2a819c3bf3df0c5814d0b3fc80ee8d",
		"note4")
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.RemoveNote(tc.txID)
			require.Equal(t, tc.expect, err)
			if err != nil {
				return
			}

			testutil.RequireFileNotExists(t, fmt.Sprintf("%s/%s.txnote", tc.manager.config.NotesDir, tc.txID))
		})
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.Equal(t, config, Config{
		NotesDir: "./notes/",
	})
}
