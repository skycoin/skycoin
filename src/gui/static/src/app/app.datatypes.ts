import { BigNumber } from 'bignumber.js';
/**
 * Internal Objects
 */

export class Address {
  address: string;
  coins: BigNumber = new BigNumber('0');
  hours: BigNumber = new BigNumber('0');
  copying?: boolean; // Optional parameter indicating whether the address is being copied to clipboard
  outputs?: any;
  confirmed?: boolean; // Optional parameter for hardware wallets only
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
  balance: BigNumber = new BigNumber('0');
  inputs: any[];
  outputs: any[];
  txid: string;
  hoursSent?: BigNumber;
  hoursBurned?: BigNumber;
  coinsMovedInternally?: boolean;
  note?: string;
}

export class PreviewTransaction extends Transaction {
  from: string;
  to: string[];
  encoded: string;
  innerHash: string;
  wallet?: Wallet;
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
  coins: BigNumber = new BigNumber('0');
  hours: BigNumber = new BigNumber('0');
  addresses: Address[];
  encrypted: boolean;
  hideEmpty?: boolean;
  opened?: boolean;
  isHardware?: boolean;
  hasHwSecurityWarnings?: boolean;
  stopShowingHwSecurityPopup?: boolean;
}

export class Connection {
  id: number;
  address: string;
  listen_port: number;
  source?: string;
}

export class TradingPair {
  from: string;
  to: string;
  price: number;
  pair: string;
  min: number;
  max: number;
}

export class ExchangeOrder {
  pair: string;
  fromAmount: number|null;
  toAmount: number;
  toAddress: string;
  toTag: string|null;
  refundAddress: string|null;
  refundTag: string|null;
  id: string;
  exchangeAddress: string;
  exchangeTag: string|null;
  toTx?: string|null;
  status: string;
  message?: string;
}

export class StoredExchangeOrder {
  id: string;
  pair: string;
  fromAmount: number;
  toAmount: number;
  address: string;
  timestamp: number;
  price: number;
}

export interface Output {
  address: string;
  coins: BigNumber;
  hash: string;
  calculated_hours: BigNumber;
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
