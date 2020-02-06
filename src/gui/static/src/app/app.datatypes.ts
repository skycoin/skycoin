/**
 * Internal Objects
 */
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
