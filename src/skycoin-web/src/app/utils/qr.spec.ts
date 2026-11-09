import { parsePaymentLink, parseScannedPayment } from './qr';

describe('parseScannedPayment', () => {
  it('reads plain addresses and payment URIs', () => {
    expect(parseScannedPayment(' 2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR ')).toEqual({ address: '2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR', amount: undefined });
    expect(parseScannedPayment('skycoin:2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR?amount=1.5&hours=2'))
      .toEqual({ address: '2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR', amount: '1.5' });
    expect(parseScannedPayment('BITCOIN:bc1qxyz?amount=.25&label=x')).toEqual({ address: 'bc1qxyz', amount: '.25' });
  });

  it('keeps only the address of an ethereum URI, whose value is in wei', () => {
    expect(parseScannedPayment('ethereum:0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC@1?value=1e18&amount=3'))
      .toEqual({ address: '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC', amount: undefined });
    expect(parseScannedPayment('ethereum:0xdAC17F958D2ee523a2206206994597C13D831ec7/transfer?address=0xabc'))
      .toEqual({ address: '0xdAC17F958D2ee523a2206206994597C13D831ec7', amount: undefined });
  });

  it('ignores an amount that is not a number', () => {
    expect(parseScannedPayment('bitcoin:bc1q?amount=1e3').amount).toBeUndefined();
  });
});

describe('parsePaymentLink', () => {
  it('reads the address, amount and hours of a payment link', () => {
    expect(parsePaymentLink(' skycoin:2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR?amount=1.5&hours=20 '))
      .toEqual({ address: '2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR', amount: '1.5', hours: '20' });
    expect(parsePaymentLink('skycoin:2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR'))
      .toEqual({ address: '2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR', amount: undefined });
  });

  it('refuses text that is not a link or has invalid hours', () => {
    expect(parsePaymentLink('2GgFvqoyk9RjwVzj8tqfcXVXB4XAd1Vb9vR')).toBeNull();
    expect(parsePaymentLink('skycoin:?amount=1')).toBeNull();
    expect(parsePaymentLink('skycoin:abc?hours=1.5')).toBeNull();
    expect(parsePaymentLink('')).toBeNull();
  });
});
