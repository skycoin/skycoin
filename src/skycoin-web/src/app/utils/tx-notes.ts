/**
 * Personal transaction notes, kept in this browser by transaction id. The node's
 * STORAGE API set is off by default and a remote node is shared, so the notes
 * stay local, as the old desktop wallet's help text promised.
 */
const storageKey = 'txNotes';

export const maxNoteChars = 64;

function readNotes(): { [txid: string]: string } {
  try {
    const notes = JSON.parse(window.localStorage.getItem(storageKey) || '{}');

    return notes && typeof notes === 'object' ? notes : {};
  } catch (e) {
    return {};
  }
}

/** The note saved for a transaction, or an empty string. */
export function getTxNote(txid: string): string {
  return (txid && readNotes()[txid]) || '';
}

/** Saves a note for a transaction, or removes it when empty. False when it could not be saved. */
export function setTxNote(txid: string, note: string): boolean {
  if (!txid) {
    return false;
  }
  try {
    const notes = readNotes();
    const text = (note || '').trim().substring(0, maxNoteChars);
    if (text) {
      notes[txid] = text;
    } else {
      delete notes[txid];
    }
    window.localStorage.setItem(storageKey, JSON.stringify(notes));

    return true;
  } catch (e) {
    return false;
  }
}

/** Sets the note property of each transaction from the saved notes. */
export function attachTxNotes(transactions: { txid: string, note?: string }[]) {
  const notes = readNotes();
  transactions.forEach(tx => tx.note = notes[tx.txid] || '');
}
