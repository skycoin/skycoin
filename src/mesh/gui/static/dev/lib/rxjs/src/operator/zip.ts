import {Observable, ObservableInput} from '../Observable';
import {ArrayObservable} from '../observable/ArrayObservable';
import {isArray} from '../util/isArray';
import {Operator} from '../Operator';
import {PartialObserver} from '../Observer';
import {Subscriber} from '../Subscriber';
import {OuterSubscriber} from '../OuterSubscriber';
import {InnerSubscriber} from '../InnerSubscriber';
import {subscribeToResult} from '../util/subscribeToResult';
import {$$iterator} from '../symbol/iterator';

/**
 * @param observables
 * @return {Observable<R>}
 * @method zip
 * @owner Observable
 */
export function zipProto<R>(...observables: Array<ObservableInput<any> | ((...values: Array<any>) => R)>): Observable<R> {
  observables.unshift(this);
  return zipStatic.apply(this, observables);
}

/* tslint:disable:max-line-length */
export interface ZipSignature<T> {
  <R>(project: (v1: T) => R): Observable<R>;
  <T2, R>(v2: ObservableInput<T2>, project: (v1: T, v2: T2) => R): Observable<R>;
  <T2, T3, R>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, project: (v1: T, v2: T2, v3: T3) => R): Observable<R>;
  <T2, T3, T4, R>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, project: (v1: T, v2: T2, v3: T3, v4: T4) => R): Observable<R>;
  <T2, T3, T4, T5, R>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, project: (v1: T, v2: T2, v3: T3, v4: T4, v5: T5) => R): Observable<R>;
  <T2, T3, T4, T5, T6, R>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, v6: ObservableInput<T6>, project: (v1: T, v2: T2, v3: T3, v4: T4, v5: T5, v6: T6) => R): Observable<R>;

  <T2>(v2: ObservableInput<T2>): Observable<[T, T2]>;
  <T2, T3>(v2: ObservableInput<T2>, v3: ObservableInput<T3>): Observable<[T, T2, T3]>;
  <T2, T3, T4>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>): Observable<[T, T2, T3, T4]>;
  <T2, T3, T4, T5>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>): Observable<[T, T2, T3, T4, T5]>;
  <T2, T3, T4, T5, T6>(v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, v6: ObservableInput<T6>): Observable<[T, T2, T3, T4, T5, T6]>;

  <R>(...observables: Array<ObservableInput<any> | ((...values: Array<any>) => R)>): Observable<R>;
  <R>(array: ObservableInput<any>[]): Observable<R>;
  <R>(array: ObservableInput<any>[], project: (...values: Array<any>) => R): Observable<R>;
}
/* tslint:enable:max-line-length */

/* tslint:disable:max-line-length */
export function zipStatic<T>(v1: ObservableInput<T>): Observable<[T]>;
export function zipStatic<T, T2>(v1: ObservableInput<T>, v2: ObservableInput<T2>): Observable<[T, T2]>;
export function zipStatic<T, T2, T3>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>): Observable<[T, T2, T3]>;
export function zipStatic<T, T2, T3, T4>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>): Observable<[T, T2, T3, T4]>;
export function zipStatic<T, T2, T3, T4, T5>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>): Observable<[T, T2, T3, T4, T5]>;
export function zipStatic<T, T2, T3, T4, T5, T6>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, v6: ObservableInput<T6>): Observable<[T, T2, T3, T4, T5, T6]>;
export function zipStatic<T, R>(v1: ObservableInput<T>, project: (v1: T) => R): Observable<R>;
export function zipStatic<T, T2, R>(v1: ObservableInput<T>, v2: ObservableInput<T2>, project: (v1: T, v2: T2) => R): Observable<R>;
export function zipStatic<T, T2, T3, R>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, project: (v1: T, v2: T2, v3: T3) => R): Observable<R>;
export function zipStatic<T, T2, T3, T4, R>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, project: (v1: T, v2: T2, v3: T3, v4: T4) => R): Observable<R>;
export function zipStatic<T, T2, T3, T4, T5, R>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, project: (v1: T, v2: T2, v3: T3, v4: T4, v5: T5) => R): Observable<R>;
export function zipStatic<T, T2, T3, T4, T5, T6, R>(v1: ObservableInput<T>, v2: ObservableInput<T2>, v3: ObservableInput<T3>, v4: ObservableInput<T4>, v5: ObservableInput<T5>, v6: ObservableInput<T6>, project: (v1: T, v2: T2, v3: T3, v4: T4, v5: T5, v6: T6) => R): Observable<R>;
export function zipStatic<R>(...observables: Array<ObservableInput<any> | ((...values: Array<any>) => R)>): Observable<R>;
export function zipStatic<R>(array: ObservableInput<any>[]): Observable<R>;
export function zipStatic<R>(array: ObservableInput<any>[], project: (...values: Array<any>) => R): Observable<R>;
/* tslint:enable:max-line-length */

/**
 * @param observables
 * @return {Observable<R>}
 * @static true
 * @name zip
 * @owner Observable
 */
export function zipStatic<T, R>(...observables: Array<ObservableInput<any> | ((...values: Array<any>) => R)>): Observable<R> {
  const project = <((...ys: Array<any>) => R)> observables[observables.length - 1];
  if (typeof project === 'function') {
    observables.pop();
  }
  return new ArrayObservable(observables).lift(new ZipOperator(project));
}

export class ZipOperator<T, R> implements Operator<T, R> {

  project: (...values: Array<any>) => R;

  constructor(project?: (...values: Array<any>) => R) {
    this.project = project;
  }

  call(subscriber: Subscriber<R>, source: any): any {
    return source._subscribe(new ZipSubscriber(subscriber, this.project));
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
export class ZipSubscriber<T, R> extends Subscriber<T> {
  private index = 0;
  private values: any;
  private project: (...values: Array<any>) => R;
  private iterators: LookAheadIterator<any>[] = [];
  private active = 0;

  constructor(destination: Subscriber<R>,
              project?: (...values: Array<any>) => R,
              values: any = Object.create(null)) {
    super(destination);
    this.project = (typeof project === 'function') ? project : null;
    this.values = values;
  }

  protected _next(value: any) {
    const iterators = this.iterators;
    const index = this.index++;
    if (isArray(value)) {
      iterators.push(new StaticArrayIterator(value));
    } else if (typeof value[$$iterator] === 'function') {
      iterators.push(new StaticIterator(value[$$iterator]()));
    } else {
      iterators.push(new ZipBufferIterator(this.destination, this, value, index));
    }
  }

  protected _complete() {
    const iterators = this.iterators;
    const len = iterators.length;
    this.active = len;
    for (let i = 0; i < len; i++) {
      let iterator: ZipBufferIterator<any, any> = <any>iterators[i];
      if (iterator.stillUnsubscribed) {
        this.add(iterator.subscribe(iterator, i));
      } else {
        this.active--; // not an observable
      }
    }
  }

  notifyInactive() {
    this.active--;
    if (this.active === 0) {
      this.destination.complete();
    }
  }

  checkIterators() {
    const iterators = this.iterators;
    const len = iterators.length;
    const destination = this.destination;

    // abort if not all of them have values
    for (let i = 0; i < len; i++) {
      let iterator = iterators[i];
      if (typeof iterator.hasValue === 'function' && !iterator.hasValue()) {
        return;
      }
    }

    let shouldComplete = false;
    const args: any[] = [];
    for (let i = 0; i < len; i++) {
      let iterator = iterators[i];
      let result = iterator.next();

      // check to see if it's completed now that you've gotten
      // the next value.
      if (iterator.hasCompleted()) {
        shouldComplete = true;
      }

      if (result.done) {
        destination.complete();
        return;
      }

      args.push(result.value);
    }

    if (this.project) {
      this._tryProject(args);
    } else {
      destination.next(args);
    }

    if (shouldComplete) {
      destination.complete();
    }
  }

  protected _tryProject(args: any[]) {
    let result: any;
    try {
      result = this.project.apply(this, args);
    } catch (err) {
      this.destination.error(err);
      return;
    }
    this.destination.next(result);
  }
}

interface LookAheadIterator<T> extends Iterator<T> {
  hasValue(): boolean;
  hasCompleted(): boolean;
}

class StaticIterator<T> implements LookAheadIterator<T> {
  private nextResult: IteratorResult<T>;

  constructor(private iterator: Iterator<T>) {
    this.nextResult = iterator.next();
  }

  hasValue() {
    return true;
  }

  next(): IteratorResult<T> {
    const result = this.nextResult;
    this.nextResult = this.iterator.next();
    return result;
  }

  hasCompleted() {
    const nextResult = this.nextResult;
    return nextResult && nextResult.done;
  }
}

class StaticArrayIterator<T> implements LookAheadIterator<T> {
  private index = 0;
  private length = 0;

  constructor(private array: T[]) {
    this.length = array.length;
  }

  [$$iterator]() {
    return this;
  }

  next(value?: any): IteratorResult<T> {
    const i = this.index++;
    const array = this.array;
    return i < this.length ? { value: array[i], done: false } : { done: true };
  }

  hasValue() {
    return this.array.length > this.index;
  }

  hasCompleted() {
    return this.array.length === this.index;
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class ZipBufferIterator<T, R> extends OuterSubscriber<T, R> implements LookAheadIterator<T> {
  stillUnsubscribed = true;
  buffer: T[] = [];
  isComplete = false;

  constructor(destination: PartialObserver<T>,
              private parent: ZipSubscriber<T, R>,
              private observable: Observable<T>,
              private index: number) {
    super(destination);
  }

  [$$iterator]() {
    return this;
  }

  // NOTE: there is actually a name collision here with Subscriber.next and Iterator.next
  //    this is legit because `next()` will never be called by a subscription in this case.
  next(): IteratorResult<T> {
    const buffer = this.buffer;
    if (buffer.length === 0 && this.isComplete) {
      return { done: true };
    } else {
      return { value: buffer.shift(), done: false };
    }
  }

  hasValue() {
    return this.buffer.length > 0;
  }

  hasCompleted() {
    return this.buffer.length === 0 && this.isComplete;
  }

  notifyComplete() {
    if (this.buffer.length > 0) {
      this.isComplete = true;
      this.parent.notifyInactive();
    } else {
      this.destination.complete();
    }
  }

  notifyNext(outerValue: T, innerValue: any,
             outerIndex: number, innerIndex: number,
             innerSub: InnerSubscriber<T, R>): void {
    this.buffer.push(innerValue);
    this.parent.checkIterators();
  }

  subscribe(value: any, index: number) {
    return subscribeToResult<any, any>(this, this.observable, this, index);
  }
}
