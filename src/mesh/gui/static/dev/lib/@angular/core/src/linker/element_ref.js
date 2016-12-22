"use strict";
/**
* A wrapper around a native element inside of a View.
*
* An `ElementRef` is backed by a render-specific element. In the browser, this is usually a DOM
* element.
*/
// Note: We don't expose things like `Injector`, `ViewContainer`, ... here,
// i.e. users have to ask for what they need. With that, we can build better analysis tools
// and could do better codegen in the future.
var ElementRef = (function () {
    function ElementRef(nativeElement) {
        this.nativeElement = nativeElement;
    }
    return ElementRef;
}());
exports.ElementRef = ElementRef;
//# sourceMappingURL=element_ref.js.map