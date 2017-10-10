import { DateFromNowPipe } from './date-from-now.pipe';

describe('DateFromNowPipe', () => {
  it('create an instance', () => {
    const pipe = new DateFromNowPipe();
    expect(pipe).toBeTruthy();
  });
});
