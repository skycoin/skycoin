import { parseBulkDestinations } from './bulk-send.component';

describe('parseBulkDestinations', () => {
  it('reads address and coins', () => {
    expect(parseBulkDestinations('a1,10\n a2 , 2.5 \n')).toEqual([
      { address: 'a1', coins: '10' },
      { address: 'a2', coins: '2.5' },
    ]);
  });

  it('reads hours when every line has them', () => {
    expect(parseBulkDestinations('a1,10,5\r\na2,1,0')).toEqual([
      { address: 'a1', coins: '10', hours: '5' },
      { address: 'a2', coins: '1', hours: '0' },
    ]);
  });

  it('skips blank lines', () => {
    expect(parseBulkDestinations('\n\na1,1\n\n')).toEqual([{ address: 'a1', coins: '1' }]);
  });

  it('rejects empty text', () => {
    expect(parseBulkDestinations('  \n ')).toBe('send.bulk-send.no-data-error');
  });

  it('rejects a line with the wrong number of fields', () => {
    expect(parseBulkDestinations('a1')).toBe('send.bulk-send.invalid-data-error');
    expect(parseBulkDestinations('a1,1,2,3')).toBe('send.bulk-send.invalid-data-error');
    expect(parseBulkDestinations('a1,,2')).toBe('send.bulk-send.invalid-data-error');
  });

  it('rejects hours on some lines only', () => {
    expect(parseBulkDestinations('a1,1,2\na2,1')).toBe('send.bulk-send.inconsistent-data-error');
  });
});
