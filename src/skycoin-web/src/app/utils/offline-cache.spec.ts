import { of, throwError } from 'rxjs';

import { withOfflineCache } from './offline-cache';

describe('withOfflineCache', () => {
  const key = 'test:spec';

  beforeEach(() => localStorage.removeItem('offlineCache:' + key));
  afterEach(() => localStorage.removeItem('offlineCache:' + key));

  it('returns a fresh response and keeps it', () => {
    let got: any;
    let at: Date | null | undefined;
    withOfflineCache(key, 'a,b', of({ n: 1 }), t => at = t).subscribe(v => got = v);
    expect(got).toEqual({ n: 1 });
    expect(at).toBeNull();
    expect(JSON.parse(localStorage.getItem('offlineCache:' + key)!).value).toEqual({ n: 1 });
  });

  it('returns the kept response with its time when the source fails', () => {
    withOfflineCache(key, 'a,b', of({ n: 1 }), () => {}).subscribe();
    let got: any;
    let at: Date | null = null;
    withOfflineCache(key, 'a,b', throwError(() => new Error('down')), t => at = t).subscribe(v => got = v);
    expect(got).toEqual({ n: 1 });
    expect(at).toEqual(jasmine.any(Date));
  });

  it('passes the error on for another scope or without a kept response', () => {
    let err: any;
    withOfflineCache(key, 'a', throwError(() => new Error('down')), () => {}).subscribe({ error: e => err = e });
    expect(err.message).toBe('down');

    withOfflineCache(key, 'a,b', of({ n: 1 }), () => {}).subscribe();
    err = null;
    withOfflineCache(key, 'a,c', throwError(() => new Error('down')), () => {}).subscribe({ error: e => err = e });
    expect(err.message).toBe('down');
  });
});
