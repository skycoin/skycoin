import { Type } from '@angular/core';
import { CanActivate } from './lifecycle_annotations_impl';
import { reflector } from '@angular/core';
export function hasLifecycleHook(e, type) {
    if (!(type instanceof Type))
        return false;
    return e.name in type.prototype;
}
export function getCanActivateHook(type) {
    var annotations = reflector.annotations(type);
    for (let i = 0; i < annotations.length; i += 1) {
        let annotation = annotations[i];
        if (annotation instanceof CanActivate) {
            return annotation.fn;
        }
    }
    return null;
}
//# sourceMappingURL=route_lifecycle_reflector.js.map