import { BaseException } from '../index';
let _FakeAsyncTestZoneSpecType = Zone['FakeAsyncTestZoneSpec'];
/**
 * Wraps a function to be executed in the fakeAsync zone:
 * - microtasks are manually executed by calling `flushMicrotasks()`,
 * - timers are synchronous, `tick()` simulates the asynchronous passage of time.
 *
 * If there are any pending timers at the end of the function, an exception will be thrown.
 *
 * Can be used to wrap inject() calls.
 *
 * ## Example
 *
 * {@example testing/ts/fake_async.ts region='basic'}
 *
 * @param fn
 * @returns {Function} The function wrapped to be executed in the fakeAsync zone
 */
export function fakeAsync(fn) {
    if (Zone.current.get('FakeAsyncTestZoneSpec') != null) {
        throw new BaseException('fakeAsync() calls can not be nested');
    }
    let fakeAsyncTestZoneSpec = new _FakeAsyncTestZoneSpecType();
    let fakeAsyncZone = Zone.current.fork(fakeAsyncTestZoneSpec);
    return function (...args) {
        let res = fakeAsyncZone.run(() => {
            let res = fn(...args);
            flushMicrotasks();
            return res;
        });
        if (fakeAsyncTestZoneSpec.pendingPeriodicTimers.length > 0) {
            throw new BaseException(`${fakeAsyncTestZoneSpec.pendingPeriodicTimers.length} ` +
                `periodic timer(s) still in the queue.`);
        }
        if (fakeAsyncTestZoneSpec.pendingTimers.length > 0) {
            throw new BaseException(`${fakeAsyncTestZoneSpec.pendingTimers.length} timer(s) still in the queue.`);
        }
        return res;
    };
}
function _getFakeAsyncZoneSpec() {
    let zoneSpec = Zone.current.get('FakeAsyncTestZoneSpec');
    if (zoneSpec == null) {
        throw new Error('The code should be running in the fakeAsync zone to call this function');
    }
    return zoneSpec;
}
/**
 * Clear the queue of pending timers and microtasks.
 * Tests no longer need to call this explicitly.
 *
 * @deprecated
 */
export function clearPendingTimers() {
    // Do nothing.
}
/**
 * Simulates the asynchronous passage of time for the timers in the fakeAsync zone.
 *
 * The microtasks queue is drained at the very start of this function and after any timer callback
 * has been executed.
 *
 * ## Example
 *
 * {@example testing/ts/fake_async.ts region='basic'}
 *
 * @param {number} millis Number of millisecond, defaults to 0
 */
export function tick(millis = 0) {
    _getFakeAsyncZoneSpec().tick(millis);
}
/**
 * Flush any pending microtasks.
 */
export function flushMicrotasks() {
    _getFakeAsyncZoneSpec().flushMicrotasks();
}
//# sourceMappingURL=fake_async.js.map