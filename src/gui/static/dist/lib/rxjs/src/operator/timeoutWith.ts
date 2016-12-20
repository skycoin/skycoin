import {Operator} from '../Operator';
import {Subscriber} from '../Subscriber';
import {Scheduler} from '../Scheduler';
import {async} from '../scheduler/async';
import {Subscription} from '../Subscription';
import {Observable} from '../Observable';
import {isDate} from '../util/isDate';
import {OuterSubscriber} from '../OuterSubscriber';
import {subscribeToResult} from '../util/subscribeToResult';

/**
 * @param due
 * @param withObservable
 * @param scheduler
 * @return {Observable<R>|WebSocketSubject<T>|Observable<T>}
 * @method timeoutWith
 * @owner Observable
 */
export function timeoutWith<T, R>(due: number | Date,
                                  withObservable: Observable<R>,
                                  scheduler: Scheduler = async): Observable<T | R> {
  let absoluteTimeout = isDate(due);
  let waitFor = absoluteTimeout ? (+due - scheduler.now()) : Math.abs(<number>due);
  return this.lift(new TimeoutWithOperator(waitFor, absoluteTimeout, withObservable, scheduler));
}

export interface TimeoutWithSignature<T> {
  (due: number | Date, withObservable: Observable<T>, scheduler?: Scheduler): Observable<T>;
  <R>(due: number | Date, withObservable: Observable<R>, scheduler?: Scheduler): Observable<T | R>;
}

class TimeoutWithOperator<T> implements Operator<T, T> {
  constructor(private waitFor: number,
              private absoluteTimeout: boolean,
              private withObservable: Observable<any>,
              private scheduler: Scheduler) {
  }

  call(subscriber: Subscriber<T>, source: any): any {
    return source._subscribe(new TimeoutWithSubscriber(
      subscriber, this.absoluteTimeout, this.waitFor, this.withObservable, this.scheduler
    ));
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class TimeoutWithSubscriber<T, R> extends OuterSubscriber<T, R> {
  private timeoutSubscription: Subscription = undefined;
  private index: number = 0;
  private _previousIndex: number = 0;
  get previousIndex(): number {
    return this._previousIndex;
  }
  private _hasCompleted: boolean = false;
  get hasCompleted(): boolean {
    return this._hasCompleted;
  }

  constructor(public destination: Subscriber<T>,
              private absoluteTimeout: boolean,
              private waitFor: number,
              private withObservable: Observable<any>,
              private scheduler: Scheduler) {
    super();
    destination.add(this);
    this.scheduleTimeout();
  }

  private static dispatchTimeout(state: any): void {
    const source = state.subscriber;
    const currentIndex = state.index;
    if (!source.hasCompleted && source.previousIndex === currentIndex) {
      source.handleTimeout();
    }
  }

  private scheduleTimeout(): void {
    let currentIndex = this.index;
    const timeoutState = { subscriber: this, index: currentIndex };
    this.scheduler.schedule(TimeoutWithSubscriber.dispatchTimeout, this.waitFor, timeoutState);
    this.index++;
    this._previousIndex = currentIndex;
  }

  protected _next(value: T) {
    this.destination.next(value);
    if (!this.absoluteTimeout) {
      this.scheduleTimeout();
    }
  }

  protected _error(err: any) {
    this.destination.error(err);
    this._hasCompleted = true;
  }

  protected _complete() {
    this.destination.complete();
    this._hasCompleted = true;
  }

  handleTimeout(): void {
    if (!this.isUnsubscribed) {
      const withObservable = this.withObservable;
      this.unsubscribe();
      this.destination.add(this.timeoutSubscription = subscribeToResult(this, withObservable));
    }
  }
}
