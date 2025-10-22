import { BigNumber } from 'bignumber.js';

export interface Wallet {
  label: string;
  addresses: Address[];
  seed?: string;
  needSeedConfirmation?: boolean;
  balance?: BigNumber;
  hours?: BigNumber;
  hidden?: boolean;
  opened?: boolean;
  hideEmpty?: boolean;
  coinId?: number;
  nextSeed?: string;
}

export interface Address {
  address: string;
  secret_key?: string;
  public_key?: string;
  balance?: BigNumber;
  hours?: BigNumber;
  outputs?: GetOutputsRequestOutput[];
  isCopying?: boolean;
}

export class Transaction {
  balance?: BigNumber;
  inputs: any[];
  outputs: any[];
  hoursSent?: BigNumber;
  hoursBurned?: BigNumber;
  coinsMovedInternally?: boolean;
}

export class NormalTransaction extends Transaction {
  txid: string;
  addresses: string[];
  timestamp: number;
  block: number;
  confirmed: boolean;
}

export class PreviewTransaction extends Transaction {
  from: string;
  to: string;
  encoded: string;
}

export interface Output {
  address: string;
  coins: BigNumber;
  hash: string;
  calculated_hours: BigNumber;
}

export interface GetOutputsRequest {
  head_outputs: GetOutputsRequestOutput[];
  outgoing_outputs: any[];
  incoming_outputs: any[];
}

export interface GetOutputsRequestOutput {
  hash: string;
  src_tx: string;
  address: string;
  coins: string;
  calculated_hours: number;
}

export class TransactionInput {
  hash: string;
  secret: string;
  address?: string;
  coins?: number;
  calculated_hours?: number;
}

export class TransactionOutput {
  address: string;
  coins: number;
  hours: number;
}

export class TotalBalance {
  coins: BigNumber = new BigNumber('0');
  hours: BigNumber = new BigNumber('0');
}

export interface Balance {
  confirmed: {
    coins: number;
    hours: number;
  };
  predicted: {
    coins: number;
    hours: number;
  };
  addresses: {
    [key: string]: AddressBalance
  };
}

export interface AddressBalance {
  confirmed: {
    coins: number;
    hours: number;
  };
  predicted: {
    coins: number;
    hours: number;
  };
}

export interface ConfirmationData {
  text: string;
  headerText: string;
  checkboxText?: string;
  confirmButtonText: string;
  cancelButtonText?: string;
  redTitle?: boolean;
  disableDismiss?: boolean;
}
