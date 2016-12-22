/**
* JS version of browser APIs. This library can only run in the browser.
*/
var win = typeof window !== 'undefined' && window || {};
export { win as window };
export var document = win.document;
export var location = win.location;
export var gc = win['gc'] ? () => win['gc']() : () => null;
export var performance = win['performance'] ? win['performance'] : null;
export const Event = win['Event'];
export const MouseEvent = win['MouseEvent'];
export const KeyboardEvent = win['KeyboardEvent'];
export const EventTarget = win['EventTarget'];
export const History = win['History'];
export const Location = win['Location'];
export const EventListener = win['EventListener'];
//# sourceMappingURL=browser.js.map