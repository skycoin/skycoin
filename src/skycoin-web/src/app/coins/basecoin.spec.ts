import { BaseCoin, skycoinArt } from './basecoin';
import { SkycoinCoin } from './skycoin.coin';

describe('BaseCoin.fromServerData', () => {
  it('gives SKY from a server the Skycoin header art', () => {
    const coin = BaseCoin.fromServerData({ id: 1, coinName: 'Skycoin', coinSymbol: 'SKY', coinType: 'skycoin' });
    expect(coin.imageName).toBe(skycoinArt.imageName);
    expect(coin.gradientName).toBe(skycoinArt.gradientName);
    expect(coin.mobileImageName).toBe(skycoinArt.mobileImageName);
  });

  it('leaves another fiber coin without it', () => {
    const coin = BaseCoin.fromServerData({ id: 2, coinName: 'Other', coinSymbol: 'OTH', coinType: 'skycoin' });
    expect(coin.imageName).toBe('');
    expect(coin.mobileImageName).toBe('');
  });

  it('matches the built-in Skycoin coin', () => {
    const builtIn = new SkycoinCoin();
    expect(builtIn.imageName).toBe(skycoinArt.imageName);
    expect(builtIn.mobileImageName).toBe(skycoinArt.mobileImageName);
  });
});
