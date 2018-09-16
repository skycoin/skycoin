package visor

import (
	"github.com/skycoin/skycoin/src/notes"
)


// GetAllNotes
func (v Visor) GetAllNotes() []notes.Note {
	return v.Notes.GetAll()
}

// If note wasn't found by Id -> return empty Note
func (v Visor) GetNoteByTransId(txId string) (notes.Note) {
	return v.Notes.GetByTransId(txId)
}

// Add Note
func (v Visor) AddNote(note notes.Note) error {
	return v.Notes.Add(note)
}

// Remove Note by TransactionId
func (v Visor) RemoveNote(txId string) error {
	return v.Notes.RemoveByTxId(txId)
}