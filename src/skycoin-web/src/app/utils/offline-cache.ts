import { Observable, of, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs';

interface CacheEntry {
  scope: string;
  time: number;
  value: any;
}

/**
 * Keeps the last good response of source under key, for one scope such as an
 * address list. When source fails, the cached response of the same scope is
 * returned instead and onResult gets its time. onResult gets null after a
 * fresh response. Without a cached response the error is passed on.
 */
export function withOfflineCache<T>(key: string, scope: string, source: Observable<T>, onResult: (cachedAt: Date | null) => void): Observable<T> {
  const storageKey = 'offlineCache:' + key;

  return source.pipe(
    tap(value => {
      onResult(null);
      try {
        const entry: CacheEntry = { scope: scope, time: Date.now(), value: value };
        localStorage.setItem(storageKey, JSON.stringify(entry));
      } catch (e) { }
    }),
    catchError(err => {
      let entry: CacheEntry | null = null;
      try {
        entry = JSON.parse(localStorage.getItem(storageKey)!);
      } catch (e) { }
      if (!entry || entry.scope !== scope) {
        return throwError(() => err);
      }
      onResult(new Date(entry.time));

      return of(entry.value as T);
    }),
  );
}
