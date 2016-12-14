/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
class ProxyZoneSpec {
    constructor(defaultSpecDelegate = null) {
        this.defaultSpecDelegate = defaultSpecDelegate;
        this.name = 'ProxyZone';
        this.properties = { 'ProxyZoneSpec': this };
        this.propertyKeys = null;
        this.setDelegate(defaultSpecDelegate);
    }
    static get() {
        return Zone.current.get('ProxyZoneSpec');
    }
    static isLoaded() {
        return ProxyZoneSpec.get() instanceof ProxyZoneSpec;
    }
    static assertPresent() {
        if (!this.isLoaded()) {
            throw new Error(`Expected to be running in 'ProxyZone', but it was not found.`);
        }
        return ProxyZoneSpec.get();
    }
    setDelegate(delegateSpec) {
        this._delegateSpec = delegateSpec;
        this.propertyKeys && this.propertyKeys.forEach((key) => delete this.properties[key]);
        this.propertyKeys = null;
        if (delegateSpec && delegateSpec.properties) {
            this.propertyKeys = Object.keys(delegateSpec.properties);
            this.propertyKeys.forEach((k) => this.properties[k] = delegateSpec.properties[k]);
        }
    }
    getDelegate() {
        return this._delegateSpec;
    }
    resetDelegate() {
        this.setDelegate(this.defaultSpecDelegate);
    }
    onFork(parentZoneDelegate, currentZone, targetZone, zoneSpec) {
        if (this._delegateSpec && this._delegateSpec.onFork) {
            return this._delegateSpec.onFork(parentZoneDelegate, currentZone, targetZone, zoneSpec);
        }
        else {
            return parentZoneDelegate.fork(targetZone, zoneSpec);
        }
    }
    onIntercept(parentZoneDelegate, currentZone, targetZone, delegate, source) {
        if (this._delegateSpec && this._delegateSpec.onIntercept) {
            return this._delegateSpec.onIntercept(parentZoneDelegate, currentZone, targetZone, delegate, source);
        }
        else {
            return parentZoneDelegate.intercept(targetZone, delegate, source);
        }
    }
    onInvoke(parentZoneDelegate, currentZone, targetZone, delegate, applyThis, applyArgs, source) {
        if (this._delegateSpec && this._delegateSpec.onInvoke) {
            return this._delegateSpec.onInvoke(parentZoneDelegate, currentZone, targetZone, delegate, applyThis, applyArgs, source);
        }
        else {
            return parentZoneDelegate.invoke(targetZone, delegate, applyThis, applyArgs, source);
        }
    }
    onHandleError(parentZoneDelegate, currentZone, targetZone, error) {
        if (this._delegateSpec && this._delegateSpec.onHandleError) {
            return this._delegateSpec.onHandleError(parentZoneDelegate, currentZone, targetZone, error);
        }
        else {
            return parentZoneDelegate.handleError(targetZone, error);
        }
    }
    onScheduleTask(parentZoneDelegate, currentZone, targetZone, task) {
        if (this._delegateSpec && this._delegateSpec.onScheduleTask) {
            return this._delegateSpec.onScheduleTask(parentZoneDelegate, currentZone, targetZone, task);
        }
        else {
            return parentZoneDelegate.scheduleTask(targetZone, task);
        }
    }
    onInvokeTask(parentZoneDelegate, currentZone, targetZone, task, applyThis, applyArgs) {
        if (this._delegateSpec && this._delegateSpec.onFork) {
            return this._delegateSpec.onInvokeTask(parentZoneDelegate, currentZone, targetZone, task, applyThis, applyArgs);
        }
        else {
            return parentZoneDelegate.invokeTask(targetZone, task, applyThis, applyArgs);
        }
    }
    onCancelTask(parentZoneDelegate, currentZone, targetZone, task) {
        if (this._delegateSpec && this._delegateSpec.onCancelTask) {
            return this._delegateSpec.onCancelTask(parentZoneDelegate, currentZone, targetZone, task);
        }
        else {
            return parentZoneDelegate.cancelTask(targetZone, task);
        }
    }
    onHasTask(delegate, current, target, hasTaskState) {
        if (this._delegateSpec && this._delegateSpec.onHasTask) {
            this._delegateSpec.onHasTask(delegate, current, target, hasTaskState);
        }
        else {
            delegate.hasTask(target, hasTaskState);
        }
    }
}
// Export the class so that new instances can be created with proper
// constructor params.
Zone['ProxyZoneSpec'] = ProxyZoneSpec;

//# sourceMappingURL=proxy.js.map
