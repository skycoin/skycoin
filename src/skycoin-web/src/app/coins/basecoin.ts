// The Skycoin header art, for the built-in coin and for SKY in a server's coin list.
export const skycoinArt = {
  imageName: 'skycoin-header.jpg',
  gradientName: 'skycoin-gradient.png',
  mobileImageName: 'skycoin-city-header.jpg',
};

// Icons of gateway coins, by coin type or by token symbol.
const gatewayIcons: Record<string, string> = {
  'bitcoin': 'btc-icon.png',
  'bitcoin-segwit': 'btc-icon.png',
  'ethereum': 'eth-icon.png',
  'USDT': 'usdt-icon.png',
};

export class BaseCoin {
  id!: number;
  nodeUrl!: string;
  coinName!: string;
  coinSymbol!: string;
  hoursName!: string;
  priceTickerId!: string;
  priceTickerSource!: string;
  coinExplorer!: string;
  coinType!: string;
  serverWallets!: boolean;
  // Servers the node settings offer for this coin, as the server lists them.
  nodeOptions: string[] = [];
  imageName!: string;
  gradientName!: string;
  // Art for the header on a phone, shown there in place of the image and gradient.
  mobileImageName = '';
  iconName!: string;
  bigIconName!: string;
  // Base units per coin as a power of ten, for Ether and tokens.
  decimals?: number;
  // An ERC-20 token's contract address.
  contract?: string;

  constructor(data?: Partial<BaseCoin>) {
    if (data) {
      Object.assign(this, data);
    }
  }

  hasHours(): boolean {
    return !this.isGateway();
  }

  isBitcoin(): boolean {
    return this.coinType === 'bitcoin' || this.coinType === 'bitcoin-segwit';
  }

  /** Ether, or an ERC-20 token held at Ethereum addresses. */
  isEthereum(): boolean {
    return this.coinType === 'ethereum' || this.coinType === 'erc20';
  }

  isToken(): boolean {
    return this.coinType === 'erc20';
  }

  /**
   * A coin skycoin web serves from its own /v1/btc/ or /v1/eth/ routes, with
   * server wallets, rather than through a Skycoin node.
   */
  isGateway(): boolean {
    return this.isBitcoin() || this.isEthereum();
  }

  /** The route prefix of a gateway coin. */
  get gatewayPrefix(): string {
    return this.isEthereum() ? 'eth' : 'btc';
  }

  get coinsMultiplier(): number {
    return Math.pow(10, this.coinDecimals);
  }

  get coinDecimals(): number {
    if (this.isEthereum()) {
      return this.decimals ?? 18;
    }

    return this.isBitcoin() ? 8 : 6;
  }

  /** Decimals worth showing; Ether has 18, more than anyone reads. */
  get displayDecimals(): number {
    return Math.min(this.coinDecimals, 8);
  }

  static fromServerData(data: any): BaseCoin {
    const coin = new BaseCoin();
    coin.id = data.id;
    coin.nodeUrl = data.nodeUrl || '';
    coin.coinName = data.coinName || 'Unknown';
    coin.coinSymbol = data.coinSymbol || '???';
    coin.hoursName = data.hoursName || 'Coin Hours';
    coin.priceTickerId = data.priceTickerId || '';
    coin.priceTickerSource = data.priceTickerSource || 'coinpaprika';
    coin.coinExplorer = data.coinExplorer || '';
    coin.coinType = data.coinType || 'skycoin';
    coin.serverWallets = !!data.serverWallets;
    coin.nodeOptions = Array.isArray(data.nodeOptions) ? data.nodeOptions.filter((s: any) => typeof s === 'string') : [];
    if (typeof data.decimals === 'number') {
      coin.decimals = data.decimals;
    }
    coin.contract = data.contract || undefined;
    // Use generic/skycoin assets as fallback for dynamically discovered coins
    coin.imageName = '';
    coin.gradientName = '';
    coin.iconName = 'skycoin-icon.png';
    coin.bigIconName = 'skycoin-icon-b.png';
    if (coin.coinType === 'skycoin' && coin.coinSymbol === 'SKY') {
      Object.assign(coin, skycoinArt);
    }
    const icon = gatewayIcons[coin.isToken() ? coin.coinSymbol : coin.coinType];
    if (icon) {
      coin.iconName = icon;
      coin.bigIconName = icon;
    }
    return coin;
  }
}
