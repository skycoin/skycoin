package visor

import (
	"github.com/skycoin/skycoin/src/notes"
)


// GetAllNotes returns all saved Notes
func (v Visor) GetAllNotes() []notes.Note {
	return v.Notes.GetAll()
}

// GetNoteByTxID If note wasn't found by Id -> return empty Note
func (v Visor) GetNoteByTxID(txID string) (notes.Note) {
	return v.Notes.GetByTxID(txID)
}

// AddNote adds a Note
func (v Visor) AddNote(note notes.Note) error {
	return v.Notes.Add(note)
}

// Remove Note by TransactionId
func (v Visor) RemoveNote(txId string) error {
	return v.Notes.RemoveByTxID(txId)
}