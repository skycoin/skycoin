import {Observable} from '../Observable';
import {Operator} from '../Operator';
import {Subscriber} from '../Subscriber';
import {EmptyError} from '../util/EmptyError';

/**
 * Returns an Observable that emits only the last item emitted by the source Observable.
 * It optionally takes a predicate function as a parameter, in which case, rather than emitting
 * the last item from the source Observable, the resulting Observable will emit the last item
 * from the source Observable that satisfies the predicate.
 *
 * <img src="./img/last.png" width="100%">
 *
 * @throws {EmptyError} Delivers an EmptyError to the Observer's `error`
 * callback if the Observable completes before any `next` notification was sent.
 * @param {function} predicate - the condition any source emitted item has to satisfy.
 * @return {Observable} an Observable that emits only the last item satisfying the given condition
 * from the source, or an NoSuchElementException if no such items are emitted.
 * @throws - Throws if no items that match the predicate are emitted by the source Observable.
 * @method last
 * @owner Observable
 */
export function last<T, R>(predicate?: (value: T, index: number, source: Observable<T>) => boolean,
                           resultSelector?: (value: T, index: number) => R | void,
                           defaultValue?: R): Observable<T | R> {
  return this.lift(new LastOperator(predicate, resultSelector, defaultValue, this));
}

export interface LastSignature<T> {
  (predicate?: (value: T, index: number, source: Observable<T>) => boolean): Observable<T>;
  (predicate: (value: T, index: number, source: Observable<T>) => boolean, resultSelector: void, defaultValue?: T): Observable<T>;
  <R>(predicate?: (value: T, index: number, source: Observable<T>) => boolean, resultSelector?: (value: T, index: number) => R,
      defaultValue?: R): Observable<R>;
}

class LastOperator<T, R> implements Operator<T, R> {
  constructor(private predicate?: (value: T, index: number, source: Observable<T>) => boolean,
              private resultSelector?: (value: T, index: number) => R,
              private defaultValue?: any,
              private source?: Observable<T>) {
  }

  call(observer: Subscriber<R>, source: any): any {
    return source._subscribe(new LastSubscriber(observer, this.predicate, this.resultSelector, this.defaultValue, this.source));
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class LastSubscriber<T, R> extends Subscriber<T> {
  private lastValue: T | R;
  private hasValue: boolean = false;
  private index: number = 0;

  constructor(destination: Subscriber<R>,
              private predicate?: (value: T, index: number, source: Observable<T>) => boolean,
              private resultSelector?: (value: T, index: number) => R,
              private defaultValue?: any,
              private source?: Observable<T>) {
    super(destination);
    if (typeof defaultValue !== 'undefined') {
      this.lastValue = defaultValue;
      this.hasValue = true;
    }
  }

  protected _next(value: T): void {
    const index = this.index++;
    if (this.predicate) {
      this._tryPredicate(value, index);
    } else {
      if (this.resultSelector) {
        this._tryResultSelector(value, index);
        return;
      }
      this.lastValue = value;
      this.hasValue = true;
    }
  }

  private _tryPredicate(value: T, index: number) {
    let result: any;
    try {
      result = this.predicate(value, index, this.source);
    } catch (err) {
      this.destination.error(err);
      return;
    }
    if (result) {
      if (this.resultSelector) {
        this._tryResultSelector(value, index);
        return;
      }
      this.lastValue = value;
      this.hasValue = true;
    }
  }

  private _tryResultSelector(value: T, index: number) {
    let result: any;
    try {
      result = this.resultSelector(value, index);
    } catch (err) {
      this.destination.error(err);
      return;
    }
    this.lastValue = result;
    this.hasValue = true;
  }

  protected _complete(): void {
    const destination = this.destination;
    if (this.hasValue) {
      destination.next(this.lastValue);
      destination.complete();
    } else {
      destination.error(new EmptyError);
    }
  }
}
