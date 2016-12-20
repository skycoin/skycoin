import { Injectable } from '@angular/core';
import { getDOM } from '../dom_adapter';
import { EventManagerPlugin } from './event_manager';
export class DomEventsPlugin extends EventManagerPlugin {
    // This plugin should come last in the list of plugins, because it accepts all
    // events.
    supports(eventName) { return true; }
    addEventListener(element, eventName, handler) {
        var zone = this.manager.getZone();
        var outsideHandler = (event) => zone.runGuarded(() => handler(event));
        return this.manager.getZone().runOutsideAngular(() => getDOM().onAndCancel(element, eventName, outsideHandler));
    }
    addGlobalEventListener(target, eventName, handler) {
        var element = getDOM().getGlobalEventTarget(target);
        var zone = this.manager.getZone();
        var outsideHandler = (event) => zone.runGuarded(() => handler(event));
        return this.manager.getZone().runOutsideAngular(() => getDOM().onAndCancel(element, eventName, outsideHandler));
    }
}
DomEventsPlugin.decorators = [
    { type: Injectable },
];
//# sourceMappingURL=dom_events.js.map