import { AmountPipe } from './amount.pipe';

describe('AmountPipe', () => {
  it('create an instance', () => {
    const pipe = new AmountPipe(null);
    expect(pipe).toBeTruthy();
  });
});
