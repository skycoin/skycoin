import { BaseCoin } from './basecoin';
import { coinsId } from '../constants/coins-id.const';

export class SkycoinCoin extends BaseCoin {
  id = coinsId.sky;
  nodeUrl = '';  // Empty - will use server proxy
  coinName = 'Skycoin';
  coinSymbol = 'SKY';
  hoursName = 'Coin Hours';
  priceTickerId = 'sky-skycoin';
  coinExplorer = 'https://explorer.skycoin.net';
  imageName = 'skycoin-header.jpg';
  gradientName = 'skycoin-gradient.png';
  iconName = 'skycoin-icon.png';
  bigIconName = 'skycoin-icon-b.png';
}
