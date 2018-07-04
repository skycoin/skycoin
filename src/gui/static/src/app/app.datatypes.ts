/**
 * Internal Objects
 */

export class Address {
  address: string;
  coins: number;
  hours: number;
  copying?: boolean; // Optional parameter indicating whether the address is being copied to clipboard
  outputs?: any;
}

export class PurchaseOrder {
  coin_type: string;
  filename: string;
  deposit_address: string;
  recipient_address: string;
  status?: string;
}

export class TellerConfig {
  enabled: boolean;
  sky_btc_exchange_rate: number;
}

export class Transaction {
  balance: number;
  inputs: any[];
  outputs: any[];
  txid: string;
  hoursSent?: number;
  hoursBurned?: number;
}

export class PreviewTransaction extends Transaction {
  from: string;
  to: string[];
  encoded: string;
}

export class NormalTransaction extends Transaction {
  addresses: string[];
  timestamp: number;
  block: number;
  confirmed: boolean;
}

export class Version {
  version: string;
}

export class Wallet {
  label: string;
  filename: string;
  coins: number;
  hours: number;
  addresses: Address[];
  encrypted: boolean;
  hideEmpty?: boolean;
  opened?: boolean;
}

export class Connection {
  id: number;
  address: string;
  listen_port: number;
  source?: string;
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
  encrypted: boolean;
}

export class GetWalletsResponseEntry {
  address: string;
}
