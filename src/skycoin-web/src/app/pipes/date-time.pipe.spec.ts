import { DateTimePipe } from './date-time.pipe';

describe('DateTimePipe', () => {
  it('create an instance', () => {
    const pipe = new DateTimePipe();
    expect(pipe).toBeTruthy();
  });
});
