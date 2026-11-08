import { Injectable } from '@angular/core';
import { Observable, forkJoin, from, of, throwError } from 'rxjs';
import { map, mergeMap } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { BlockchainService } from '../blockchain.service';
import { CipherProvider } from '../cipher.provider';
import { Output, Wallet } from '../../app.datatypes';
import { WalletService, walletAddresses } from './wallet.service';

/** An output of a joint transaction, with coins and hours set by hand. */
export interface JointDestination {
  address: string;
  coins: BigNumber;
  hours: BigNumber;
}

export interface JointInput {
  index: number;
  hash: string;
  address: string;
  coins: BigNumber;
  hours: BigNumber;
  signed: boolean;
  mine: boolean;
}

export interface JointOutput {
  address: string;
  coins: BigNumber;
  hours: BigNumber;
  mine: boolean;
}

/** A decoded joint transaction, as the node sees it. */
export interface JointTransaction {
  encoded: string;
  txid: string;
  inputs: JointInput[];
  outputs: JointOutput[];
  inputHours: BigNumber;
  fee: BigNumber;
  minimumFee: BigNumber;
  signedCount: number;
  fullySigned: boolean;
}

/** The parts of an encoded transaction that signing needs. */
export interface RawTransaction {
  innerHash: string;
  inputs: string[];
  signed: boolean[];
}

// The layout of a serialized coin.Transaction (src/coin/transactions.go): length
// u32, type u8, inner hash, then u32 counted signatures, inputs and outputs.
const sigsCountOffset = 4 + 1 + 32;
const sigLength = 65;
const hashLength = 32;

function readUint32(bytes: Uint8Array, offset: number): number {
  if (offset + 4 > bytes.length) {
    throw new Error('joint-tx.invalid-tx');
  }

  return (bytes[offset] | (bytes[offset + 1] << 8) | (bytes[offset + 2] << 16)) + bytes[offset + 3] * 0x1000000;
}

function hexToBytes(hex: string): Uint8Array {
  if (!/^([0-9a-fA-F]{2})*$/.test(hex)) {
    throw new Error('joint-tx.invalid-tx');
  }
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
  }

  return bytes;
}

function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
}

/** Reads the inner hash, inputs and signing state of an encoded classic transaction. */
export function parseRawTransaction(encoded: string): RawTransaction {
  const bytes = hexToBytes(encoded.trim());
  if (bytes.length < sigsCountOffset + 4 || readUint32(bytes, 0) !== bytes.length || bytes[4] !== 0) {
    throw new Error('joint-tx.invalid-tx');
  }

  const nSigs = readUint32(bytes, sigsCountOffset);
  const signed: boolean[] = [];
  let offset = sigsCountOffset + 4;
  for (let i = 0; i < nSigs; i++, offset += sigLength) {
    if (offset + sigLength > bytes.length) {
      throw new Error('joint-tx.invalid-tx');
    }
    signed.push(bytes.subarray(offset, offset + sigLength).some(b => b !== 0));
  }

  const nIn = readUint32(bytes, offset);
  offset += 4;
  if (nIn !== nSigs || offset + nIn * hashLength > bytes.length) {
    throw new Error('joint-tx.invalid-tx');
  }
  const inputs: string[] = [];
  for (let i = 0; i < nIn; i++, offset += hashLength) {
    inputs.push(bytesToHex(bytes.subarray(offset, offset + hashLength)));
  }

  return { innerHash: bytesToHex(bytes.subarray(5, 5 + hashLength)), inputs, signed };
}

/** Returns encoded with signature sig at input index. */
export function setSignature(encoded: string, index: number, sig: string): string {
  if (!/^[0-9a-fA-F]{130}$/.test(sig)) {
    throw new Error('joint-tx.invalid-signature');
  }
  const start = (sigsCountOffset + 4 + index * sigLength) * 2;
  const tx = encoded.trim();

  return tx.substring(0, start) + sig.toLowerCase() + tx.substring(start + sigLength * 2);
}

/** The least fee the node accepts for inputs holding hours: hours / burnFactor, rounded up. */
export function minimumFee(hours: BigNumber, burnFactor: BigNumber): BigNumber {
  if (burnFactor.isLessThanOrEqualTo(0)) {
    return new BigNumber(0);
  }

  return hours.dividedBy(burnFactor).integerValue(BigNumber.ROUND_CEIL);
}

/** Splits pasted output hashes on commas and white space. */
export function parseHashes(text: string): string[] {
  return (text || '').split(/[\s,]+/).map(h => h.trim().toLowerCase()).filter(h => h.length > 0);
}

/**
 * Builds, reviews and co-signs joint transactions: one transaction spending
 * outputs of several owners, each of whom signs only their own inputs.
 */
@Injectable()
export class JointTransactionService {

  constructor(
    private apiService: ApiService,
    private blockchainService: BlockchainService,
    private cipherProvider: CipherProvider,
    private walletService: WalletService,
  ) {}

  /** The unspent outputs with these hashes, in no particular order. */
  outputs(hashes: string[]): Observable<Output[]> {
    if (hashes.length === 0) {
      return of([]);
    }

    return this.apiService.post('outputs', { hashes: hashes.join(',') }).pipe(map((response: any) =>
      (response.head_outputs || []).map((o: any) => ({
        address: o.address,
        coins: new BigNumber(o.coins),
        hash: o.hash,
        calculated_hours: new BigNumber(o.calculated_hours),
      }))));
  }

  minimumFee(hours: BigNumber): BigNumber {
    return minimumFee(hours, this.blockchainService.burnRate);
  }

  /** Creates the unsigned transaction spending every one of unspents. Returns it encoded. */
  create(unspents: string[], destinations: JointDestination[], changeAddress: string | null): Observable<string> {
    const params: any = {
      hours_selection: { type: 'manual' },
      unspents: unspents,
      to: destinations.map(d => ({ address: d.address, coins: d.coins.toString(), hours: d.hours.toFixed(0) })),
    };
    if (changeAddress) {
      params.change_address = changeAddress;
    }

    return this.apiService.post('transaction', params, { json: true }, true).pipe(mergeMap((response: any) => {
      const used = new Set<string>(response.data.transaction.inputs.map((i: any) => i.uxid));
      if (unspents.some(u => !used.has(u))) {
        return throwError(() => new Error('joint-tx.input-left-out'));
      }

      return of(response.data.encoded_transaction as string);
    }));
  }

  /** Decodes encoded and checks it on the node, marking what wallet owns. */
  decode(encoded: string, wallet: Wallet | null): Observable<JointTransaction> {
    let raw: RawTransaction;
    try {
      raw = parseRawTransaction(encoded);
    } catch (e) {
      return throwError(() => e);
    }
    const fullySigned = raw.signed.length > 0 && raw.signed.every(s => s);
    const body = { encoded_transaction: encoded.trim(), unsigned: !fullySigned };
    const mine = new Set<string>(wallet ? walletAddresses(wallet).map(a => a.address) : []);

    return this.apiService.post('transaction/verify', body, { json: true }, true).pipe(map((response: any) => {
      const tx = response.data.transaction;
      const inputs: JointInput[] = tx.inputs.map((i: any, index: number) => ({
        index,
        hash: i.uxid,
        address: i.address,
        coins: new BigNumber(i.coins),
        hours: new BigNumber(i.calculated_hours),
        signed: raw.signed[index],
        mine: mine.has(i.address),
      }));
      const outputs: JointOutput[] = tx.outputs.map((o: any) => ({
        address: o.address,
        coins: new BigNumber(o.coins),
        hours: new BigNumber(o.hours),
        mine: mine.has(o.address),
      }));
      const inputHours = inputs.reduce((sum, i) => sum.plus(i.hours), new BigNumber(0));

      return {
        encoded: body.encoded_transaction,
        txid: tx.txid,
        inputs,
        outputs,
        inputHours,
        fee: new BigNumber(tx.fee),
        minimumFee: this.minimumFee(inputHours),
        signedCount: raw.signed.filter(s => s).length,
        fullySigned,
      };
    }));
  }

  /** Signs the unsigned inputs of tx that wallet owns. Returns the encoded result. */
  sign(wallet: Wallet, password: string | undefined, tx: JointTransaction): Observable<string> {
    const toSign = tx.inputs.filter(i => i.mine && !i.signed);
    if (toSign.length === 0) {
      return throwError(() => new Error('joint-tx.nothing-to-sign'));
    }

    if (this.walletService.isServerWallet(wallet)) {
      const body = {
        wallet_id: wallet.filename,
        encoded_transaction: tx.encoded,
        input_addresses: tx.inputs.map(i => i.address),
        sign_indexes: toSign.map(i => i.index),
        password: password || '',
      };

      return this.apiService.post('wallet/transaction/sign', body, { json: true }, true).pipe(
        map((response: any) => response.data.encoded_transaction as string));
    }

    const addresses = walletAddresses(wallet);
    const innerHash = parseRawTransaction(tx.encoded).innerHash;
    const signatures = toSign.map(input => {
      const address = addresses.find(a => a.address === input.address);
      if (!address || !address.secret_key) {
        return throwError(() => new Error('offline-tx.not-in-wallet'));
      }

      return from(inputHash(innerHash, input.hash)).pipe(
        mergeMap(hash => this.cipherProvider.signHash(hash, address.secret_key!)),
        map(sig => ({ index: input.index, sig })));
    });

    return forkJoin(signatures).pipe(map(sigs =>
      sigs.reduce((encoded, s) => setSignature(encoded, s.index, s.sig), tx.encoded)));
  }
}

/** The hash an input's signature covers: SHA256(inner hash || input hash), as cipher.AddSHA256. */
export async function inputHash(innerHash: string, input: string): Promise<string> {
  if (!crypto.subtle) {
    throw new Error('errors.insecure-context');
  }
  const data = new Uint8Array(hashLength * 2);
  data.set(hexToBytes(innerHash), 0);
  data.set(hexToBytes(input), hashLength);
  const digest = await crypto.subtle.digest('SHA-256', data as BufferSource);

  return bytesToHex(new Uint8Array(digest));
}
