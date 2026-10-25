import { BaseCoin, skycoinArt } from './basecoin';
import { defaultCoinId } from '../constants/coins-id.const';

/**
 * Fallback Skycoin coin definition used when the server is not available.
 */
export class SkycoinCoin extends BaseCoin {
  constructor() {
    super({
      id: defaultCoinId,
      // Empty means the node that served the wallet, as with the node's -web-wallet.
      nodeUrl: '',
      coinName: 'Skycoin',
      coinSymbol: 'SKY',
      hoursName: 'Coin Hours',
      priceTickerId: 'sky-skycoin',
      priceTickerSource: 'coinpaprika',
      coinExplorer: 'https://explorer.skycoin.net',
      coinType: 'skycoin',
      ...skycoinArt,
      iconName: 'skycoin-icon.png',
      bigIconName: 'skycoin-icon-b.png',
    });
  }
}
