import { of, throwError, lastValueFrom } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { CipherProvider } from '../cipher.provider';
import { loadCipherWasm } from '../../utils/wasm-test-utils';
import {
  JointTransactionService, inputHash, minimumFee, parseHashes, parseRawTransaction, setSignature,
} from './joint-transaction.service';

// A transaction with two inputs and input 1 signed, from src/coin.
const unsigned =
  '3d01000000ba1edd7c833942af86bd5cdf4137aaf56ea8c512304db21ca067848942ba88190200000000000000000000000000000000000000000000' +
  '000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' +
  '000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000001ca4b3296f' +
  '0b31b9f584e06c4dfc59119ecc408ba0395047c8d616f4c88a20c687dac51506d06652be02110bd1e34c2156faa5fa37253dc885da93b840b4bdec02' +
  '00000000057a62c64df07d9525dec462acddea6690a602ece803000000000000000000000000000000057a62c64df07d9525dec462acddea6690a602' +
  'ecd0070000000000000500000000000000';
const partial =
  '3d01000000ba1edd7c833942af86bd5cdf4137aaf56ea8c512304db21ca067848942ba88190200000000000000000000000000000000000000000000' +
  '00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000eefb69e214aab8e5abe7efd7c61b' +
  '72508ee1db7caafb5275b66a5b1831589ad944b93d9638fd21794f5c62cf346fc9326f4ae70381ce1e87cfa792112e29f07a01020000001ca4b3296f' +
  '0b31b9f584e06c4dfc59119ecc408ba0395047c8d616f4c88a20c687dac51506d06652be02110bd1e34c2156faa5fa37253dc885da93b840b4bdec02' +
  '00000000057a62c64df07d9525dec462acddea6690a602ece803000000000000000000000000000000057a62c64df07d9525dec462acddea6690a602' +
  'ecd0070000000000000500000000000000';
const innerHash = 'ba1edd7c833942af86bd5cdf4137aaf56ea8c512304db21ca067848942ba8819';
const in0 = '1ca4b3296f0b31b9f584e06c4dfc59119ecc408ba0395047c8d616f4c88a20c6';
const in1 = '87dac51506d06652be02110bd1e34c2156faa5fa37253dc885da93b840b4bdec';
const hash1 = 'f2d1a5806f10ee5071017201babacd96bfcbfbd09fabce0597d83371eb3691e4';
const sig1 = 'eefb69e214aab8e5abe7efd7c61b72508ee1db7caafb5275b66a5b1831589ad944b93d9638fd21794f5c62cf346fc9326f4ae70381ce1e87cfa792112e29f07a01';
const secret = 'f7e13071cfc5707a00b1e4d40922e717c02caafc35f9f491e235aa0f669e812d';
const mine = '3Cr8gvoZ6yVQdidTPHTQ89yQ5dtBpqnm8v';
const other = '2YnwXdZNvREm5EPVkpt1u3hNuNuAUao7xb';

describe('joint transaction encoding', () => {
  it('reads the inner hash, inputs and signatures', () => {
    expect(parseRawTransaction(unsigned)).toEqual({ innerHash, inputs: [in0, in1], signed: [false, false] });
    expect(parseRawTransaction(partial).signed).toEqual([false, true]);
  });

  it('rejects what is not a classic transaction', () => {
    expect(() => parseRawTransaction('zz')).toThrowError('joint-tx.invalid-tx');
    expect(() => parseRawTransaction(unsigned.slice(0, 100))).toThrowError('joint-tx.invalid-tx');
    expect(() => parseRawTransaction(unsigned + '00')).toThrowError('joint-tx.invalid-tx');
  });

  it('writes a signature in place', () => {
    expect(setSignature(unsigned, 1, sig1)).toBe(partial);
    expect(() => setSignature(unsigned, 1, 'abc')).toThrowError('joint-tx.invalid-signature');
  });

  it('hashes an input as cipher.AddSHA256', async () => {
    expect(await inputHash(innerHash, in1)).toBe(hash1);
  });

  it('signs an input with the wasm cipher so the signature checks out', async () => {
    await loadCipherWasm();
    const hash = await inputHash(innerHash, in1);
    const sig = await lastValueFrom(new CipherProvider(null as any).signHash(hash, secret));
    expect((window as any).SkycoinCipherExtras.verifyAddressSignedHash(mine, sig, hash)).toBeNull();
    expect(parseRawTransaction(setSignature(unsigned, 1, sig)).signed).toEqual([false, true]);
  }, 60000);

  it('rounds the minimum fee up', () => {
    expect(minimumFee(new BigNumber(1117), new BigNumber(10)).toString()).toBe('112');
    expect(minimumFee(new BigNumber(1110), new BigNumber(10)).toString()).toBe('111');
  });

  it('splits pasted hashes', () => {
    expect(parseHashes(' AB, cd\nef  ')).toEqual(['ab', 'cd', 'ef']);
    expect(parseHashes('')).toEqual([]);
  });
});

describe('JointTransactionService', () => {
  let api: jasmine.SpyObj<any>;
  let cipher: jasmine.SpyObj<any>;
  let serverWallet: boolean;
  let service: JointTransactionService;

  const verifyResponse = (signed: boolean) => ({
    data: {
      unsigned: !signed,
      transaction: {
        txid: 'txid',
        fee: '117',
        inputs: [
          { uxid: in0, address: other, coins: '0.001000', calculated_hours: '5' },
          { uxid: in1, address: mine, coins: '0.001000', calculated_hours: '1112' },
        ],
        outputs: [
          { address: mine, coins: '0.001000', hours: '0' },
          { address: other, coins: '0.001000', hours: '1000' },
        ],
      },
    },
  });

  beforeEach(() => {
    api = jasmine.createSpyObj('ApiService', ['post']);
    cipher = jasmine.createSpyObj('CipherProvider', ['signHash']);
    serverWallet = false;
    service = new JointTransactionService(
      api,
      { burnRate: new BigNumber(10) } as any,
      cipher,
      { isServerWallet: () => serverWallet } as any,
    );
  });

  it('creates through the node with manual hours and a change address', async () => {
    api.post.and.returnValue(of({ data: { encoded_transaction: 'hex', transaction: { inputs: [{ uxid: in0 }, { uxid: in1 }] } } }));
    const encoded = await lastValueFrom(service.create([in0, in1],
      [{ address: other, coins: new BigNumber('0.001'), hours: new BigNumber(1000) }], mine));
    expect(encoded).toBe('hex');
    expect(api.post).toHaveBeenCalledWith('transaction', {
      hours_selection: { type: 'manual' },
      unspents: [in0, in1],
      to: [{ address: other, coins: '0.001', hours: '1000' }],
      change_address: mine,
    }, { json: true }, true);
  });

  it('fails when the node leaves an input out', async () => {
    api.post.and.returnValue(of({ data: { encoded_transaction: 'hex', transaction: { inputs: [{ uxid: in0 }] } } }));
    await expectAsync(lastValueFrom(service.create([in0, in1], [], null))).toBeRejectedWithError('joint-tx.input-left-out');
  });

  it('decodes and marks the wallet inputs', async () => {
    api.post.and.returnValue(of(verifyResponse(false)));
    const tx = await lastValueFrom(service.decode(partial, { label: 'w', addresses: [{ address: mine }] }));
    expect(api.post).toHaveBeenCalledWith('transaction/verify', { encoded_transaction: partial, unsigned: true }, { json: true }, true);
    expect(tx.inputs.map(i => [i.mine, i.signed])).toEqual([[false, false], [true, true]]);
    expect(tx.outputs.map(o => o.mine)).toEqual([true, false]);
    expect(tx.inputHours.toString()).toBe('1117');
    expect(tx.minimumFee.toString()).toBe('112');
    expect(tx.signedCount).toBe(1);
    expect(tx.fullySigned).toBeFalse();
  });

  it('signs only its own inputs of a browser wallet', async () => {
    api.post.and.returnValue(of(verifyResponse(false)));
    cipher.signHash.and.returnValue(of(sig1));
    const wallet = { label: 'w', addresses: [{ address: mine, secret_key: secret }] };
    const tx = await lastValueFrom(service.decode(unsigned, wallet));
    const signed = await lastValueFrom(service.sign(wallet, undefined, tx));
    expect(cipher.signHash).toHaveBeenCalledOnceWith(hash1, secret);
    expect(signed).toBe(partial);
  });

  it('refuses when the wallet owns no unsigned input', async () => {
    api.post.and.returnValue(of(verifyResponse(false)));
    const tx = await lastValueFrom(service.decode(partial, { label: 'w', addresses: [{ address: mine }] }));
    await expectAsync(lastValueFrom(service.sign({ label: 'w', addresses: [] }, undefined, tx)))
      .toBeRejectedWithError('joint-tx.nothing-to-sign');
  });

  it('asks a server wallet to sign by index', async () => {
    serverWallet = true;
    api.post.and.returnValues(of(verifyResponse(false)), of({ data: { encoded_transaction: 'signed' } }));
    const wallet = { label: 'w', filename: 'w.wlt', addresses: [{ address: mine }] };
    const tx = await lastValueFrom(service.decode(unsigned, wallet));
    expect(await lastValueFrom(service.sign(wallet, 'pw', tx))).toBe('signed');
    expect(api.post.calls.mostRecent().args).toEqual(['wallet/transaction/sign', {
      wallet_id: 'w.wlt',
      encoded_transaction: unsigned,
      input_addresses: [other, mine],
      sign_indexes: [1],
      password: 'pw',
    }, { json: true }, true]);
  });

  it('passes on what the node says about a bad transaction', async () => {
    api.post.and.returnValue(throwError(() => new Error('Invalid signature')));
    await expectAsync(lastValueFrom(service.decode(partial, null))).toBeRejectedWithError('Invalid signature');
  });
});
