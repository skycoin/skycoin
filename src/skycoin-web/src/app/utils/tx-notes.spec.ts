import { attachTxNotes, getTxNote, maxNoteChars, setTxNote } from './tx-notes';

describe('tx notes', () => {
  let saved: string | null;

  beforeEach(() => {
    saved = window.localStorage.getItem('txNotes');
    window.localStorage.removeItem('txNotes');
  });

  afterEach(() => {
    if (saved === null) {
      window.localStorage.removeItem('txNotes');
    } else {
      window.localStorage.setItem('txNotes', saved);
    }
  });

  it('saves, reads and removes a note', () => {
    expect(getTxNote('a')).toBe('');
    expect(setTxNote('a', '  rent  ')).toBe(true);
    expect(getTxNote('a')).toBe('rent');
    expect(setTxNote('a', ' ')).toBe(true);
    expect(getTxNote('a')).toBe('');
  });

  it('cuts a note to the maximum length', () => {
    setTxNote('a', 'x'.repeat(maxNoteChars + 10));
    expect(getTxNote('a').length).toBe(maxNoteChars);
  });

  it('attaches notes to transactions', () => {
    setTxNote('a', 'one');
    const txs: { txid: string, note?: string }[] = [{ txid: 'a' }, { txid: 'b' }];
    attachTxNotes(txs);
    expect(txs[0].note).toBe('one');
    expect(txs[1].note).toBe('');
  });

  it('survives a broken stored value', () => {
    window.localStorage.setItem('txNotes', '{not json');
    expect(getTxNote('a')).toBe('');
    expect(setTxNote('a', 'ok')).toBe(true);
    expect(getTxNote('a')).toBe('ok');
  });
});
