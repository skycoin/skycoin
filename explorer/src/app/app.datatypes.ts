/* eslint-disable @typescript-eslint/naming-convention */

import { BigNumber } from 'bignumber.js';

/**
 * Collection of types used in various parts of the application. Some types are just
 * representations of the responses the API return and others are for internal use in the
 * application. Also, this file contains functions for processing the responses returned
 * by the API and return internal types.
 */

/**
 * Elementary types used in the app.
 */

export class Block {
  id: number;
  hash: string;
  parent_hash: string;
  timestamp: number;
  transactions: Transaction[];
  size: number;
}

export class Blockchain {
  blocks: number;
}

export class Input {
  owner: string;
  coins: BigNumber;
  uxid: string;
  hours: BigNumber;
  calculatedHours: BigNumber;
}

export class Output {
  address: string;
  coins: BigNumber;
  hash: string;
  hours: BigNumber;
}

export class Transaction {
  block: number;
  id: string;
  inputs: Input[];
  outputs: Output[];
  status: boolean;
  timestamp: number;
  balance: BigNumber;
  initialBalance: BigNumber;
  finalBalance: BigNumber;
  length: number;
  fee: BigNumber;
}

export class RichlistEntry {
  address: string;
  coins: string;
  locked: boolean;
}

/**
 * Generic API response types
 */

export class GenericBlockResponse {
  header: GenericBlockHeaderResponse;
  body: GenericBlockBodyResponse;
  size: number;
}

class GenericBlockHeaderResponse {
  block_hash: string;
  previous_block_hash: string;
  seq: number;
  timestamp: number;
}

class GenericBlockBodyResponse {
  txns: GenericTransactionResponse[];
}

export class GenericTransactionResponse {
  inputs: GenericTransactionInputResponse[];
  outputs: GenericTransactionOutputResponse[];
  status: any;
  timestamp: number;
  txid: string;
  length: number;
  fee: number;
}

class GenericTransactionInputResponse {
  uxid: string;
  owner: string;
  coins: string;
  hours: number;
  calculated_hours: number;
}

class GenericTransactionOutputResponse {
  coins: string;
  dst: string;
  hours: number;
  uxid: string;
}

export function parseGenericBlock(block: GenericBlockResponse): Block {
  return {
    id: block.header.seq,
    hash: block.header.block_hash,
    parent_hash: block.header.previous_block_hash,
    timestamp: block.header.timestamp,
    transactions: block.body.txns.map(transaction => parseGenericTransaction(transaction)),
    size: block.size,
  };
}

export function parseGenericTransaction(raw: GenericTransactionResponse, address: string = null): Transaction {
  let balance = null;
  if (address) {
    balance = new BigNumber('0');
    for (const input of raw.inputs) {
      if (input.owner.toLowerCase() === address.toLowerCase()) {
        balance = balance.minus(input.coins);
      }
    }
    for (const output of raw.outputs) {
      if (output.dst.toLowerCase() === address.toLowerCase()) {
        balance = balance.plus(output.coins);
      }
    }
  }

  const response = {
    block: null,
    id: raw.txid,
    timestamp: raw.timestamp,
    inputs: raw.inputs.map(input => parseGenericTransactionInput(input)),
    outputs: raw.outputs.map(output => parseGenericTransactionOutput(output)),
    status: null,
    balance: balance,
    initialBalance: null,
    finalBalance: null,
    length: raw.length,
    fee: new BigNumber(raw.fee),
  };

  if (raw.status) {
    if (raw.status.confirmed) {
      response.status = raw.status.confirmed;
    }

    if (raw.status.height) {
      response.block = raw.status.block_seq;
    }
  }

  return response;
}

function parseGenericTransactionInput(raw: GenericTransactionInputResponse): Input {
  return {
    owner: raw.owner,
    coins: new BigNumber(raw.coins),
    uxid: raw.uxid,
    hours: new BigNumber(raw.hours),
    calculatedHours: new BigNumber(raw.calculated_hours),
  };
}

function parseGenericTransactionOutput(raw: GenericTransactionOutputResponse): Output {
  return {
    address: raw.dst,
    coins: new BigNumber(raw.coins),
    hash: raw.uxid,
    hours: new BigNumber(raw.hours),
  };
}

/**
 * API response types (returned by a specific API endpoint)
 */

export class GetUnconfirmedTransactionResponse {
  transaction: GenericTransactionResponse;
  received: string;
  is_valid: boolean;
}

export function parseGetUnconfirmedTransaction(raw: GetUnconfirmedTransactionResponse): Transaction {
  raw.transaction.timestamp = new Date(raw.received).getTime();
  raw.transaction.status = { confirmed: raw.is_valid };

  return parseGenericTransaction(raw.transaction);
}

export class GetBlocksResponse {
  blocks: GenericBlockResponse[];
}

export class GetBlockchainMetadataResponse {
  head: GetBlockchainMetadataResponseHead;
}

class GetBlockchainMetadataResponseHead {
  seq: number;
}

export class GetBalanceResponse {
  confirmed: GetBalanceResponseElement;
  predicted: GetBalanceResponseElement;
}

class GetBalanceResponseElement {
  coins: number;
  hours: number;
}

export class GetCurrentBalanceResponse {
  head_outputs: GetCurrentBalanceResponseOutput[];
}

class GetCurrentBalanceResponseOutput {
  hash: string;
  src_tx: string;
  address: string;
  coins: string;
  hours: number;
  calculated_hours: number;
}

export class GetTransactionResponse {
  status: any;
  time: number;
  txn: GenericTransactionResponse;
}

export class GetSyncStateResponse {
  current: number;
  highest: number;
}

export function parseGetTransaction(raw: GetTransactionResponse, address: string = null): Transaction {
  raw.txn.status = raw.status;
  return parseGenericTransaction(raw.txn, address);
}

export class GetUxoutResponse {
  coins: number;
  hours: number;
  owner_address: string;
  uxid: string;
}

export function parseGetUxout(raw: GetUxoutResponse): Output {
  return {
    address: raw.owner_address,
    coins: new BigNumber(raw.coins).dividedBy(1000000),
    hours: new BigNumber(raw.hours),
    hash: raw.uxid,
  };
}

export interface GetCoinSupplyResponse {
  current_supply: number;
  total_supply: number;
  max_supply: number;
  current_coinhour_supply: number;
  total_coinhour_supply: number;
}

/**
 * Specific objects.
 */

export interface AddressTransactionsResponse {
  addressHasManyTransactions: boolean;
  totalTransactionsCount: number;
  currentPageIndex: number;
  totalPages: number;
  recoveredTransactions: Transaction[];
  originalTransactions: any;
}
