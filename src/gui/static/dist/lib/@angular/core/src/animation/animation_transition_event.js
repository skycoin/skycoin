/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
/**
 * An instance of this class is returned as an event parameter when an animation
 * callback is captured for an animation either during the start or done phase.
 *
 * ```typescript
 * @Component({
 *   host: {
 *     '[@myAnimationTrigger]': 'someExpression',
 *     '(@myAnimationTrigger.start)': 'captureStartEvent($event)',
 *     '(@myAnimationTrigger.done)': 'captureDoneEvent($event)',
 *   },
 *   animations: [
 *     trigger("myAnimationTrigger", [
 *        // ...
 *     ])
 *   ]
 * })
 * class MyComponent {
 *   someExpression: any = false;
 *   captureStartEvent(event: AnimationTransitionEvent) {
 *     // the toState, fromState and totalTime data is accessible from the event variable
 *   }
 *
 *   captureDoneEvent(event: AnimationTransitionEvent) {
 *     // the toState, fromState and totalTime data is accessible from the event variable
 *   }
 * }
 * ```
 *
 * @experimental Animation support is experimental.
 */
export var AnimationTransitionEvent = (function () {
    function AnimationTransitionEvent(_a) {
        var fromState = _a.fromState, toState = _a.toState, totalTime = _a.totalTime, phaseName = _a.phaseName;
        this.fromState = fromState;
        this.toState = toState;
        this.totalTime = totalTime;
        this.phaseName = phaseName;
    }
    return AnimationTransitionEvent;
}());
//# sourceMappingURL=animation_transition_event.js.map