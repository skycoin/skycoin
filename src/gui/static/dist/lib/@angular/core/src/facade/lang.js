"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var globalScope;
if (typeof window === 'undefined') {
    if (typeof WorkerGlobalScope !== 'undefined' && self instanceof WorkerGlobalScope) {
        // TODO: Replace any with WorkerGlobalScope from lib.webworker.d.ts #3492
        globalScope = self;
    }
    else {
        globalScope = global;
    }
}
else {
    globalScope = window;
}
function scheduleMicroTask(fn) {
    Zone.current.scheduleMicroTask('scheduleMicrotask', fn);
}
exports.scheduleMicroTask = scheduleMicroTask;
exports.IS_DART = false;
// Need to declare a new variable for global here since TypeScript
// exports the original value of the symbol.
var _global = globalScope;
exports.global = _global;
exports.Type = Function;
function getTypeNameForDebugging(type) {
    if (type['name']) {
        return type['name'];
    }
    return typeof type;
}
exports.getTypeNameForDebugging = getTypeNameForDebugging;
exports.Math = _global.Math;
exports.Date = _global.Date;
var _devMode = true;
var _modeLocked = false;
function lockMode() {
    _modeLocked = true;
}
exports.lockMode = lockMode;
/**
 * Disable Angular's development mode, which turns off assertions and other
 * checks within the framework.
 *
 * One important assertion this disables verifies that a change detection pass
 * does not result in additional changes to any bindings (also known as
 * unidirectional data flow).
 */
function enableProdMode() {
    if (_modeLocked) {
        // Cannot use BaseException as that ends up importing from facade/lang.
        throw 'Cannot enable prod mode after platform setup.';
    }
    _devMode = false;
}
exports.enableProdMode = enableProdMode;
function assertionsEnabled() {
    return _devMode;
}
exports.assertionsEnabled = assertionsEnabled;
// TODO: remove calls to assert in production environment
// Note: Can't just export this and import in in other files
// as `assert` is a reserved keyword in Dart
_global.assert = function assert(condition) {
    // TODO: to be fixed properly via #2830, noop for now
};
function isPresent(obj) {
    return obj !== undefined && obj !== null;
}
exports.isPresent = isPresent;
function isBlank(obj) {
    return obj === undefined || obj === null;
}
exports.isBlank = isBlank;
function isBoolean(obj) {
    return typeof obj === "boolean";
}
exports.isBoolean = isBoolean;
function isNumber(obj) {
    return typeof obj === "number";
}
exports.isNumber = isNumber;
function isString(obj) {
    return typeof obj === "string";
}
exports.isString = isString;
function isFunction(obj) {
    return typeof obj === "function";
}
exports.isFunction = isFunction;
function isType(obj) {
    return isFunction(obj);
}
exports.isType = isType;
function isStringMap(obj) {
    return typeof obj === 'object' && obj !== null;
}
exports.isStringMap = isStringMap;
var STRING_MAP_PROTO = Object.getPrototypeOf({});
function isStrictStringMap(obj) {
    return isStringMap(obj) && Object.getPrototypeOf(obj) === STRING_MAP_PROTO;
}
exports.isStrictStringMap = isStrictStringMap;
function isPromise(obj) {
    return obj instanceof _global.Promise;
}
exports.isPromise = isPromise;
function isArray(obj) {
    return Array.isArray(obj);
}
exports.isArray = isArray;
function isDate(obj) {
    return obj instanceof exports.Date && !isNaN(obj.valueOf());
}
exports.isDate = isDate;
function noop() { }
exports.noop = noop;
function stringify(token) {
    if (typeof token === 'string') {
        return token;
    }
    if (token === undefined || token === null) {
        return '' + token;
    }
    if (token.name) {
        return token.name;
    }
    if (token.overriddenName) {
        return token.overriddenName;
    }
    var res = token.toString();
    var newLineIndex = res.indexOf("\n");
    return (newLineIndex === -1) ? res : res.substring(0, newLineIndex);
}
exports.stringify = stringify;
// serialize / deserialize enum exist only for consistency with dart API
// enums in typescript don't need to be serialized
function serializeEnum(val) {
    return val;
}
exports.serializeEnum = serializeEnum;
function deserializeEnum(val, values) {
    return val;
}
exports.deserializeEnum = deserializeEnum;
function resolveEnumToken(enumValue, val) {
    return enumValue[val];
}
exports.resolveEnumToken = resolveEnumToken;
var StringWrapper = (function () {
    function StringWrapper() {
    }
    StringWrapper.fromCharCode = function (code) { return String.fromCharCode(code); };
    StringWrapper.charCodeAt = function (s, index) { return s.charCodeAt(index); };
    StringWrapper.split = function (s, regExp) { return s.split(regExp); };
    StringWrapper.equals = function (s, s2) { return s === s2; };
    StringWrapper.stripLeft = function (s, charVal) {
        if (s && s.length) {
            var pos = 0;
            for (var i = 0; i < s.length; i++) {
                if (s[i] != charVal)
                    break;
                pos++;
            }
            s = s.substring(pos);
        }
        return s;
    };
    StringWrapper.stripRight = function (s, charVal) {
        if (s && s.length) {
            var pos = s.length;
            for (var i = s.length - 1; i >= 0; i--) {
                if (s[i] != charVal)
                    break;
                pos--;
            }
            s = s.substring(0, pos);
        }
        return s;
    };
    StringWrapper.replace = function (s, from, replace) {
        return s.replace(from, replace);
    };
    StringWrapper.replaceAll = function (s, from, replace) {
        return s.replace(from, replace);
    };
    StringWrapper.slice = function (s, from, to) {
        if (from === void 0) { from = 0; }
        if (to === void 0) { to = null; }
        return s.slice(from, to === null ? undefined : to);
    };
    StringWrapper.replaceAllMapped = function (s, from, cb) {
        return s.replace(from, function () {
            var matches = [];
            for (var _i = 0; _i < arguments.length; _i++) {
                matches[_i - 0] = arguments[_i];
            }
            // Remove offset & string from the result array
            matches.splice(-2, 2);
            // The callback receives match, p1, ..., pn
            return cb(matches);
        });
    };
    StringWrapper.contains = function (s, substr) { return s.indexOf(substr) != -1; };
    StringWrapper.compare = function (a, b) {
        if (a < b) {
            return -1;
        }
        else if (a > b) {
            return 1;
        }
        else {
            return 0;
        }
    };
    return StringWrapper;
}());
exports.StringWrapper = StringWrapper;
var StringJoiner = (function () {
    function StringJoiner(parts) {
        if (parts === void 0) { parts = []; }
        this.parts = parts;
    }
    StringJoiner.prototype.add = function (part) { this.parts.push(part); };
    StringJoiner.prototype.toString = function () { return this.parts.join(""); };
    return StringJoiner;
}());
exports.StringJoiner = StringJoiner;
var NumberParseError = (function (_super) {
    __extends(NumberParseError, _super);
    function NumberParseError(message) {
        _super.call(this);
        this.message = message;
    }
    NumberParseError.prototype.toString = function () { return this.message; };
    return NumberParseError;
}(Error));
exports.NumberParseError = NumberParseError;
var NumberWrapper = (function () {
    function NumberWrapper() {
    }
    NumberWrapper.toFixed = function (n, fractionDigits) { return n.toFixed(fractionDigits); };
    NumberWrapper.equal = function (a, b) { return a === b; };
    NumberWrapper.parseIntAutoRadix = function (text) {
        var result = parseInt(text);
        if (isNaN(result)) {
            throw new NumberParseError("Invalid integer literal when parsing " + text);
        }
        return result;
    };
    NumberWrapper.parseInt = function (text, radix) {
        if (radix == 10) {
            if (/^(\-|\+)?[0-9]+$/.test(text)) {
                return parseInt(text, radix);
            }
        }
        else if (radix == 16) {
            if (/^(\-|\+)?[0-9ABCDEFabcdef]+$/.test(text)) {
                return parseInt(text, radix);
            }
        }
        else {
            var result = parseInt(text, radix);
            if (!isNaN(result)) {
                return result;
            }
        }
        throw new NumberParseError("Invalid integer literal when parsing " + text + " in base " +
            radix);
    };
    // TODO: NaN is a valid literal but is returned by parseFloat to indicate an error.
    NumberWrapper.parseFloat = function (text) { return parseFloat(text); };
    Object.defineProperty(NumberWrapper, "NaN", {
        get: function () { return NaN; },
        enumerable: true,
        configurable: true
    });
    NumberWrapper.isNaN = function (value) { return isNaN(value); };
    NumberWrapper.isInteger = function (value) { return Number.isInteger(value); };
    return NumberWrapper;
}());
exports.NumberWrapper = NumberWrapper;
exports.RegExp = _global.RegExp;
var RegExpWrapper = (function () {
    function RegExpWrapper() {
    }
    RegExpWrapper.create = function (regExpStr, flags) {
        if (flags === void 0) { flags = ''; }
        flags = flags.replace(/g/g, '');
        return new _global.RegExp(regExpStr, flags + 'g');
    };
    RegExpWrapper.firstMatch = function (regExp, input) {
        // Reset multimatch regex state
        regExp.lastIndex = 0;
        return regExp.exec(input);
    };
    RegExpWrapper.test = function (regExp, input) {
        regExp.lastIndex = 0;
        return regExp.test(input);
    };
    RegExpWrapper.matcher = function (regExp, input) {
        // Reset regex state for the case
        // someone did not loop over all matches
        // last time.
        regExp.lastIndex = 0;
        return { re: regExp, input: input };
    };
    RegExpWrapper.replaceAll = function (regExp, input, replace) {
        var c = regExp.exec(input);
        var res = '';
        regExp.lastIndex = 0;
        var prev = 0;
        while (c) {
            res += input.substring(prev, c.index);
            res += replace(c);
            prev = c.index + c[0].length;
            regExp.lastIndex = prev;
            c = regExp.exec(input);
        }
        res += input.substring(prev);
        return res;
    };
    return RegExpWrapper;
}());
exports.RegExpWrapper = RegExpWrapper;
var RegExpMatcherWrapper = (function () {
    function RegExpMatcherWrapper() {
    }
    RegExpMatcherWrapper.next = function (matcher) {
        return matcher.re.exec(matcher.input);
    };
    return RegExpMatcherWrapper;
}());
exports.RegExpMatcherWrapper = RegExpMatcherWrapper;
var FunctionWrapper = (function () {
    function FunctionWrapper() {
    }
    FunctionWrapper.apply = function (fn, posArgs) { return fn.apply(null, posArgs); };
    return FunctionWrapper;
}());
exports.FunctionWrapper = FunctionWrapper;
// JS has NaN !== NaN
function looseIdentical(a, b) {
    return a === b || typeof a === "number" && typeof b === "number" && isNaN(a) && isNaN(b);
}
exports.looseIdentical = looseIdentical;
// JS considers NaN is the same as NaN for map Key (while NaN !== NaN otherwise)
// see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map
function getMapKey(value) {
    return value;
}
exports.getMapKey = getMapKey;
function normalizeBlank(obj) {
    return isBlank(obj) ? null : obj;
}
exports.normalizeBlank = normalizeBlank;
function normalizeBool(obj) {
    return isBlank(obj) ? false : obj;
}
exports.normalizeBool = normalizeBool;
function isJsObject(o) {
    return o !== null && (typeof o === "function" || typeof o === "object");
}
exports.isJsObject = isJsObject;
function print(obj) {
    console.log(obj);
}
exports.print = print;
function warn(obj) {
    console.warn(obj);
}
exports.warn = warn;
// Can't be all uppercase as our transpiler would think it is a special directive...
var Json = (function () {
    function Json() {
    }
    Json.parse = function (s) { return _global.JSON.parse(s); };
    Json.stringify = function (data) {
        // Dart doesn't take 3 arguments
        return _global.JSON.stringify(data, null, 2);
    };
    return Json;
}());
exports.Json = Json;
var DateWrapper = (function () {
    function DateWrapper() {
    }
    DateWrapper.create = function (year, month, day, hour, minutes, seconds, milliseconds) {
        if (month === void 0) { month = 1; }
        if (day === void 0) { day = 1; }
        if (hour === void 0) { hour = 0; }
        if (minutes === void 0) { minutes = 0; }
        if (seconds === void 0) { seconds = 0; }
        if (milliseconds === void 0) { milliseconds = 0; }
        return new exports.Date(year, month - 1, day, hour, minutes, seconds, milliseconds);
    };
    DateWrapper.fromISOString = function (str) { return new exports.Date(str); };
    DateWrapper.fromMillis = function (ms) { return new exports.Date(ms); };
    DateWrapper.toMillis = function (date) { return date.getTime(); };
    DateWrapper.now = function () { return new exports.Date(); };
    DateWrapper.toJson = function (date) { return date.toJSON(); };
    return DateWrapper;
}());
exports.DateWrapper = DateWrapper;
function setValueOnPath(global, path, value) {
    var parts = path.split('.');
    var obj = global;
    while (parts.length > 1) {
        var name = parts.shift();
        if (obj.hasOwnProperty(name) && isPresent(obj[name])) {
            obj = obj[name];
        }
        else {
            obj = obj[name] = {};
        }
    }
    if (obj === undefined || obj === null) {
        obj = {};
    }
    obj[parts.shift()] = value;
}
exports.setValueOnPath = setValueOnPath;
var _symbolIterator = null;
function getSymbolIterator() {
    if (isBlank(_symbolIterator)) {
        if (isPresent(globalScope.Symbol) && isPresent(Symbol.iterator)) {
            _symbolIterator = Symbol.iterator;
        }
        else {
            // es6-shim specific logic
            var keys = Object.getOwnPropertyNames(Map.prototype);
            for (var i = 0; i < keys.length; ++i) {
                var key = keys[i];
                if (key !== 'entries' && key !== 'size' &&
                    Map.prototype[key] === Map.prototype['entries']) {
                    _symbolIterator = key;
                }
            }
        }
    }
    return _symbolIterator;
}
exports.getSymbolIterator = getSymbolIterator;
function evalExpression(sourceUrl, expr, declarations, vars) {
    var fnBody = declarations + "\nreturn " + expr + "\n//# sourceURL=" + sourceUrl;
    var fnArgNames = [];
    var fnArgValues = [];
    for (var argName in vars) {
        fnArgNames.push(argName);
        fnArgValues.push(vars[argName]);
    }
    return new (Function.bind.apply(Function, [void 0].concat(fnArgNames.concat(fnBody))))().apply(void 0, fnArgValues);
}
exports.evalExpression = evalExpression;
function isPrimitive(obj) {
    return !isJsObject(obj);
}
exports.isPrimitive = isPrimitive;
function hasConstructor(value, type) {
    return value.constructor === type;
}
exports.hasConstructor = hasConstructor;
function bitWiseOr(values) {
    return values.reduce(function (a, b) { return a | b; });
}
exports.bitWiseOr = bitWiseOr;
function bitWiseAnd(values) {
    return values.reduce(function (a, b) { return a & b; });
}
exports.bitWiseAnd = bitWiseAnd;
function escape(s) {
    return _global.encodeURI(s);
}
exports.escape = escape;
//# sourceMappingURL=lang.js.map