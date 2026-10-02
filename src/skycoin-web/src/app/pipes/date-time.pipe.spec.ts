import { DateTimePipe } from './date-time.pipe';

describe('DateTimePipe', () => {
  const pipe = new DateTimePipe();

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('formats a unix time as local YYYY-MM-DD HH:mm', () => {
    const local = new Date(2026, 0, 5, 7, 9, 30);
    expect(pipe.transform(local.getTime() / 1000)).toBe('2026-01-05 07:09');
    expect(pipe.transform(String(Math.floor(local.getTime() / 1000)))).toBe('2026-01-05 07:09');
  });

  it('reports an unusable value the way it always has', () => {
    expect(pipe.transform(undefined)).toBe('Invalid date');
    expect(pipe.transform('not a time')).toBe('Invalid date');
  });
});
