/**
 * Internal Objects
 */

export class Wallet {
  label: string;
  filename: string;
  seed: string;
  coins: number;
  hours: number;
  addresses: Address[];
  visible?: boolean;
  hideEmpty?: boolean;
}

export class Address {
  address: string;
  coins: number;
  hours: number;
}

export class Transaction {
  addresses: string[];
  balance: number;
  block: number;
  confirmed: boolean;
  inputs: any[];
  outputs: any[];
  timestamp: number;
  txid: string;
}

/**
 * Response Objects
 */

export class GetWalletsResponseWallet {
  meta: GetWalletsResponseMeta;
  entries: GetWalletsResponseEntry[];
}

export class PostWalletNewAddressResponse {
  addresses: string[];
}

/**
 * Response Embedded Objects
 */

export class GetWalletsResponseMeta {
  label: string;
  filename: string;
  seed: string;
}

export class GetWalletsResponseEntry {
  address: string;
}
