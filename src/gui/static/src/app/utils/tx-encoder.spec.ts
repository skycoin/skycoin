import { readJSON } from 'karma-read-json';
import { TxEncoder } from './tx-encoder';

fdescribe('TxEncoder', () => {

  describe('check encoding', () => {
    const txs = readJSON('test-fixtures/encoded-txs.json').txs;

    for (let i = 0; i < txs.length; i++) {
      it('encode tx ' + i, () => {
        expect(TxEncoder.encode(txs[i].inputs, txs[i].outputs, txs[i].signatures, txs[i].innerHash)).toBe(txs[i].raw);
      });
    }
  });
});
