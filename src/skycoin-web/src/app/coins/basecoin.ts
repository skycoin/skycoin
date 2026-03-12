export class BaseCoin {
  id: number;
  nodeUrl: string;
  coinName: string;
  coinSymbol: string;
  hoursName: string;
  priceTickerId: string;
  priceTickerSource: string;
  coinExplorer: string;
  imageName: string;
  gradientName: string;
  iconName: string;
  bigIconName: string;

  constructor(data?: Partial<BaseCoin>) {
    if (data) {
      Object.assign(this, data);
    }
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
    coin.imageName = 'default-header.jpg';
    coin.gradientName = 'default-gradient.png';
    coin.iconName = 'default-icon.png';
    coin.bigIconName = 'default-icon-b.png';
    return coin;
  }
}
