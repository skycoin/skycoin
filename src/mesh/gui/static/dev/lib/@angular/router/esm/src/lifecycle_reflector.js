import { Type, isBlank } from './facade/lang';
export function hasLifecycleHook(name, obj) {
    if (isBlank(obj))
        return false;
    let type = obj.constructor;
    if (!(type instanceof Type))
        return false;
    return name in type.prototype;
}
//# sourceMappingURL=lifecycle_reflector.js.map