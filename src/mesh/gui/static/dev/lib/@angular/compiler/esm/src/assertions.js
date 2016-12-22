import { isArray, isString, isBlank, assertionsEnabled } from '../src/facade/lang';
import { BaseException } from '../src/facade/exceptions';
export function assertArrayOfStrings(identifier, value) {
    if (!assertionsEnabled() || isBlank(value)) {
        return;
    }
    if (!isArray(value)) {
        throw new BaseException(`Expected '${identifier}' to be an array of strings.`);
    }
    for (var i = 0; i < value.length; i += 1) {
        if (!isString(value[i])) {
            throw new BaseException(`Expected '${identifier}' to be an array of strings.`);
        }
    }
}
//# sourceMappingURL=assertions.js.map