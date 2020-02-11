import { processErrorMsg } from './errors';

describe('errors', () => {
  const message1 = '400 Bad Request - error description.';
  const message2 = '403 Forbidden - error description.';
  const message3 = '500 Internal Server Error - error description.';

  it('parses message from 400 and 403 responses', () => {
    expect(processErrorMsg(message1)).toEqual('Error description.');
    expect(processErrorMsg(message2)).toEqual('Error description.');
  });

  it('does not parse message from other responses', () => {
    expect(processErrorMsg(message3)).toEqual(message3);
  });
});
