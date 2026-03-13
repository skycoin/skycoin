import { TxEncoder } from './tx-encoder';
import BigNumber from 'bignumber.js';

import encodedTxsFixture from '../../../test-fixtures/encoded-txs.json';

describe('TxEncoder', () => {

  describe('check encoding', () => {
    const txs = encodedTxsFixture.txs;

    for (let i = 0; i < txs.length; i++) {
      it('encode tx ' + i, () => {
        (txs[i].outputs as any[]).forEach(output => {
          output.coins = new BigNumber(output.coins).dividedBy(1000000).toString();
          output.hours = new BigNumber(output.hours).toString();
        });

        expect(TxEncoder.encode(txs[i].inputs as any, txs[i].outputs as any, txs[i].signatures, txs[i].innerHash)).toBe(txs[i].raw);
      });
    }
  });
});
