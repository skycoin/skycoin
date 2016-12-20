import { Observable } from '../Observable';
import { Operator } from '../Operator';
import { Subscriber } from '../Subscriber';
import { OuterSubscriber } from '../OuterSubscriber';
import { InnerSubscriber } from '../InnerSubscriber';
/**
 * Returns an Observable that mirrors the first source Observable to emit an item
 * from the combination of this Observable and supplied Observables
 * @param {...Observables} ...observables sources used to race for which Observable emits first.
 * @return {Observable} an Observable that mirrors the output of the first Observable to emit an item.
 * @method race
 * @owner Observable
 */
export declare function race<T>(...observables: Array<Observable<T> | Array<Observable<T>>>): Observable<T>;
export interface RaceSignature<T> {
    (...observables: Array<Observable<T> | Array<Observable<T>>>): Observable<T>;
    <R>(...observables: Array<Observable<any> | Array<Observable<T>>>): Observable<R>;
}
/**
 * Returns an Observable that mirrors the first source Observable to emit an item.
 * @param {...Observables} ...observables sources used to race for which Observable emits first.
 * @return {Observable} an Observable that mirrors the output of the first Observable to emit an item.
 * @static true
 * @name race
 * @owner Observable
 */
export declare function raceStatic<T>(...observables: Array<Observable<T> | Array<Observable<T>>>): Observable<T>;
export declare class RaceOperator<T> implements Operator<T, T> {
    call(subscriber: Subscriber<T>, source: any): any;
}
/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
export declare class RaceSubscriber<T> extends OuterSubscriber<T, T> {
    private hasFirst;
    private observables;
    private subscriptions;
    constructor(destination: Subscriber<T>);
    protected _next(observable: any): void;
    protected _complete(): void;
    notifyNext(outerValue: T, innerValue: T, outerIndex: number, innerIndex: number, innerSub: InnerSubscriber<T, T>): void;
}
