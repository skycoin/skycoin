import { BaseCoin } from './basecoin';

/**
 * Fallback Bitcoin coin definition used when the server provides Bitcoin support.
 */
export class BitcoinCoin extends BaseCoin {
  constructor() {
    super({
      id: -2,
      nodeUrl: '',
      coinName: 'Bitcoin',
      coinSymbol: 'BTC',
      hoursName: '',
      priceTickerId: 'btc-bitcoin',
      priceTickerSource: 'coinpaprika',
      coinExplorer: 'https://blockchair.com/bitcoin',
      coinType: 'bitcoin',
      imageName: '',
      gradientName: '',
      iconName: 'skycoin-icon.png',
      bigIconName: 'skycoin-icon-b.png',
    });
  }
}
