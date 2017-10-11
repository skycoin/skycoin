import { TellerStatusPipe } from './teller-status.pipe';

describe('TellerStatusPipe', () => {
  it('create an instance', () => {
    const pipe = new TellerStatusPipe();
    expect(pipe).toBeTruthy();
  });
});
