import { parseResponseMessage } from './errors';

fdescribe('errors', () => {
  const message1 = '400 Bad Request - error description';
  const message2 = '403 Forbidden - error description';
  const message3 = '500 Internal Server Error - error description';

  it('parses message from 400 and 403 responses', () => {
    expect(parseResponseMessage(message1)).toEqual('Error description');
    expect(parseResponseMessage(message2)).toEqual('Error description');
  });

  it('does not parse message from other responses', () => {
    expect(parseResponseMessage(message3)).toEqual(message3);
  });
});
