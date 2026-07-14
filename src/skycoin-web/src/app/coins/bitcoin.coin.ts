import { BaseCoin } from './basecoin';

/**
 * Fallback Bitcoin coin definition used when the server provides Bitcoin support.
 */
export class BitcoinCoin extends BaseCoin {
  constructor() {
    super({
      id: -2,
      // For bitcoin the "node URL" is the ssl:// electrum server the visor's BTC
      // gateway connects to (it reads the server straight from the request
      // origin). A public default so BTC works out of the box; change it in
      // Settings → Node.
      nodeUrl: 'ssl://electrum.blockstream.info:50002',
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
