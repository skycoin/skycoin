import {Operator} from '../Operator';
import {Subscriber} from '../Subscriber';
import {Observable} from '../Observable';

/**
 * Returns an Observable that skips `n` items emitted by an Observable.
 *
 * <img src="./img/skip.png" width="100%">
 *
 * @param {Number} the `n` of times, items emitted by source Observable should be skipped.
 * @return {Observable} an Observable that skips values emitted by the source Observable.
 *
 * @method skip
 * @owner Observable
 */
export function skip<T>(total: number): Observable<T> {
  return this.lift(new SkipOperator(total));
}

export interface SkipSignature<T> {
  (total: number): Observable<T>;
}

class SkipOperator<T> implements Operator<T, T> {
  constructor(private total: number) {
  }

  call(subscriber: Subscriber<T>, source: any): any {
    return source._subscribe(new SkipSubscriber(subscriber, this.total));
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class SkipSubscriber<T> extends Subscriber<T> {
  count: number = 0;

  constructor(destination: Subscriber<T>, private total: number) {
    super(destination);
  }

  protected _next(x: T) {
    if (++this.count > this.total) {
      this.destination.next(x);
    }
  }
}
