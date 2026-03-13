import { BaseCoin } from './basecoin';

export class TestCoin extends BaseCoin {
  constructor() {
    super({
      id: -1,
      nodeUrl: '',
      coinName: 'Testcoin',
      coinSymbol: 'TEST',
      hoursName: 'Test Hours',
      priceTickerId: 'btc-bitcoin',
      priceTickerSource: 'coinpaprika',
      coinExplorer: 'https://explorer.testcoin.net',
      imageName: 'testcoin-header.jpg',
      gradientName: 'testcoin-gradient.png',
      iconName: 'testcoin-icon.png',
      bigIconName: 'testcoin-icon-b.png',
    });
  }
}
