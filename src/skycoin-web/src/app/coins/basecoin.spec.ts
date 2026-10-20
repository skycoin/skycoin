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

describe('BaseCoin gateway coins', () => {
  const eth = BaseCoin.fromServerData({ id: 1, coinSymbol: 'ETH', coinType: 'ethereum', decimals: 18 });
  const usdt = BaseCoin.fromServerData({ id: 2, coinSymbol: 'USDT', coinType: 'erc20', decimals: 6, contract: '0xdAC17F958D2ee523a2206206994597C13D831ec7' });
  const btc = BaseCoin.fromServerData({ id: 3, coinSymbol: 'BTC', coinType: 'bitcoin' });
  const sky = BaseCoin.fromServerData({ id: 0, coinSymbol: 'SKY', coinType: 'skycoin' });

  it('tells gateway coins from Skycoin', () => {
    expect([eth, usdt, btc].every(c => c.isGateway() && !c.hasHours())).toBeTrue();
    expect(sky.isGateway()).toBeFalse();
    expect(sky.hasHours()).toBeTrue();
    expect(eth.gatewayPrefix).toBe('eth');
    expect(usdt.gatewayPrefix).toBe('eth');
    expect(btc.gatewayPrefix).toBe('btc');
    expect(usdt.isToken()).toBeTrue();
    expect(usdt.contract).toBe('0xdAC17F958D2ee523a2206206994597C13D831ec7');
  });

  it('uses each coin\'s decimals', () => {
    expect(eth.coinsMultiplier).toBe(1e18);
    expect(eth.displayDecimals).toBe(8);
    expect(usdt.coinsMultiplier).toBe(1e6);
    expect(btc.coinsMultiplier).toBe(1e8);
    expect(sky.coinsMultiplier).toBe(1e6);
  });

  it('gives each gateway coin its icon', () => {
    expect(eth.iconName).toBe('eth-icon.png');
    expect(usdt.iconName).toBe('usdt-icon.png');
    expect(btc.iconName).toBe('btc-icon.png');
  });
});
