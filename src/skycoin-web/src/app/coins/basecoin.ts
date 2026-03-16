export class BaseCoin {
  id: number;
  nodeUrl: string;
  coinName: string;
  coinSymbol: string;
  hoursName: string;
  priceTickerId: string;
  priceTickerSource: string;
  coinExplorer: string;
  coinType: string;
  imageName: string;
  gradientName: string;
  iconName: string;
  bigIconName: string;

  constructor(data?: Partial<BaseCoin>) {
    if (data) {
      Object.assign(this, data);
    }
  }

  hasHours(): boolean {
    return this.coinType !== 'bitcoin';
  }

  isBitcoin(): boolean {
    return this.coinType === 'bitcoin';
  }

  get coinsMultiplier(): number {
    return this.coinType === 'bitcoin' ? 100000000 : 1000000;
  }

  get coinDecimals(): number {
    return this.coinType === 'bitcoin' ? 8 : 6;
  }

  static fromServerData(data: any): BaseCoin {
    const isBtc = data.coinType === 'bitcoin';
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
    // Use generic/skycoin assets as fallback for dynamically discovered coins
    coin.imageName = '';
    coin.gradientName = '';
    coin.iconName = 'skycoin-icon.png';
    coin.bigIconName = 'skycoin-icon-b.png';
    return coin;
  }
}
