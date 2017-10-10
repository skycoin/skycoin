import { TransactionsAmountPipe } from './transactions-amount.pipe';

describe('TransactionsAmountPipe', () => {
  it('create an instance', () => {
    const pipe = new TransactionsAmountPipe();
    expect(pipe).toBeTruthy();
  });
});
