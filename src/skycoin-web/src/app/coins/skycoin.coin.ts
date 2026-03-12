import { BaseCoin } from './basecoin';
import { defaultCoinId } from '../constants/coins-id.const';

/**
 * Fallback Skycoin coin definition used when the server is not available.
 */
export class SkycoinCoin extends BaseCoin {
  constructor() {
    super({
      id: defaultCoinId,
      nodeUrl: '',
      coinName: 'Skycoin',
      coinSymbol: 'SKY',
      hoursName: 'Coin Hours',
      priceTickerId: 'sky-skycoin',
      priceTickerSource: 'coinpaprika',
      coinExplorer: 'https://explorer.skycoin.net',
      imageName: 'skycoin-header.jpg',
      gradientName: 'skycoin-gradient.png',
      iconName: 'skycoin-icon.png',
      bigIconName: 'skycoin-icon-b.png',
    });
  }
}
