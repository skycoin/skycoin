"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var lang_1 = require('../src/facade/lang');
var collection_1 = require('../src/facade/collection');
var parse_util_1 = require('./parse_util');
var html_tags_1 = require('./html_tags');
(function (HtmlTokenType) {
    HtmlTokenType[HtmlTokenType["TAG_OPEN_START"] = 0] = "TAG_OPEN_START";
    HtmlTokenType[HtmlTokenType["TAG_OPEN_END"] = 1] = "TAG_OPEN_END";
    HtmlTokenType[HtmlTokenType["TAG_OPEN_END_VOID"] = 2] = "TAG_OPEN_END_VOID";
    HtmlTokenType[HtmlTokenType["TAG_CLOSE"] = 3] = "TAG_CLOSE";
    HtmlTokenType[HtmlTokenType["TEXT"] = 4] = "TEXT";
    HtmlTokenType[HtmlTokenType["ESCAPABLE_RAW_TEXT"] = 5] = "ESCAPABLE_RAW_TEXT";
    HtmlTokenType[HtmlTokenType["RAW_TEXT"] = 6] = "RAW_TEXT";
    HtmlTokenType[HtmlTokenType["COMMENT_START"] = 7] = "COMMENT_START";
    HtmlTokenType[HtmlTokenType["COMMENT_END"] = 8] = "COMMENT_END";
    HtmlTokenType[HtmlTokenType["CDATA_START"] = 9] = "CDATA_START";
    HtmlTokenType[HtmlTokenType["CDATA_END"] = 10] = "CDATA_END";
    HtmlTokenType[HtmlTokenType["ATTR_NAME"] = 11] = "ATTR_NAME";
    HtmlTokenType[HtmlTokenType["ATTR_VALUE"] = 12] = "ATTR_VALUE";
    HtmlTokenType[HtmlTokenType["DOC_TYPE"] = 13] = "DOC_TYPE";
    HtmlTokenType[HtmlTokenType["EXPANSION_FORM_START"] = 14] = "EXPANSION_FORM_START";
    HtmlTokenType[HtmlTokenType["EXPANSION_CASE_VALUE"] = 15] = "EXPANSION_CASE_VALUE";
    HtmlTokenType[HtmlTokenType["EXPANSION_CASE_EXP_START"] = 16] = "EXPANSION_CASE_EXP_START";
    HtmlTokenType[HtmlTokenType["EXPANSION_CASE_EXP_END"] = 17] = "EXPANSION_CASE_EXP_END";
    HtmlTokenType[HtmlTokenType["EXPANSION_FORM_END"] = 18] = "EXPANSION_FORM_END";
    HtmlTokenType[HtmlTokenType["EOF"] = 19] = "EOF";
})(exports.HtmlTokenType || (exports.HtmlTokenType = {}));
var HtmlTokenType = exports.HtmlTokenType;
var HtmlToken = (function () {
    function HtmlToken(type, parts, sourceSpan) {
        this.type = type;
        this.parts = parts;
        this.sourceSpan = sourceSpan;
    }
    return HtmlToken;
}());
exports.HtmlToken = HtmlToken;
var HtmlTokenError = (function (_super) {
    __extends(HtmlTokenError, _super);
    function HtmlTokenError(errorMsg, tokenType, span) {
        _super.call(this, span, errorMsg);
        this.tokenType = tokenType;
    }
    return HtmlTokenError;
}(parse_util_1.ParseError));
exports.HtmlTokenError = HtmlTokenError;
var HtmlTokenizeResult = (function () {
    function HtmlTokenizeResult(tokens, errors) {
        this.tokens = tokens;
        this.errors = errors;
    }
    return HtmlTokenizeResult;
}());
exports.HtmlTokenizeResult = HtmlTokenizeResult;
function tokenizeHtml(sourceContent, sourceUrl, tokenizeExpansionForms) {
    if (tokenizeExpansionForms === void 0) { tokenizeExpansionForms = false; }
    return new _HtmlTokenizer(new parse_util_1.ParseSourceFile(sourceContent, sourceUrl), tokenizeExpansionForms)
        .tokenize();
}
exports.tokenizeHtml = tokenizeHtml;
var $EOF = 0;
var $TAB = 9;
var $LF = 10;
var $FF = 12;
var $CR = 13;
var $SPACE = 32;
var $BANG = 33;
var $DQ = 34;
var $HASH = 35;
var $$ = 36;
var $AMPERSAND = 38;
var $SQ = 39;
var $MINUS = 45;
var $SLASH = 47;
var $0 = 48;
var $SEMICOLON = 59;
var $9 = 57;
var $COLON = 58;
var $LT = 60;
var $EQ = 61;
var $GT = 62;
var $QUESTION = 63;
var $LBRACKET = 91;
var $RBRACKET = 93;
var $LBRACE = 123;
var $RBRACE = 125;
var $COMMA = 44;
var $A = 65;
var $F = 70;
var $X = 88;
var $Z = 90;
var $a = 97;
var $f = 102;
var $z = 122;
var $x = 120;
var $NBSP = 160;
var CR_OR_CRLF_REGEXP = /\r\n?/g;
function unexpectedCharacterErrorMsg(charCode) {
    var char = charCode === $EOF ? 'EOF' : lang_1.StringWrapper.fromCharCode(charCode);
    return "Unexpected character \"" + char + "\"";
}
function unknownEntityErrorMsg(entitySrc) {
    return "Unknown entity \"" + entitySrc + "\" - use the \"&#<decimal>;\" or  \"&#x<hex>;\" syntax";
}
var ControlFlowError = (function () {
    function ControlFlowError(error) {
        this.error = error;
    }
    return ControlFlowError;
}());
// See http://www.w3.org/TR/html51/syntax.html#writing
var _HtmlTokenizer = (function () {
    function _HtmlTokenizer(file, tokenizeExpansionForms) {
        this.file = file;
        this.tokenizeExpansionForms = tokenizeExpansionForms;
        // Note: this is always lowercase!
        this.peek = -1;
        this.nextPeek = -1;
        this.index = -1;
        this.line = 0;
        this.column = -1;
        this.expansionCaseStack = [];
        this.tokens = [];
        this.errors = [];
        this.input = file.content;
        this.length = file.content.length;
        this._advance();
    }
    _HtmlTokenizer.prototype._processCarriageReturns = function (content) {
        // http://www.w3.org/TR/html5/syntax.html#preprocessing-the-input-stream
        // In order to keep the original position in the source, we can not
        // pre-process it.
        // Instead CRs are processed right before instantiating the tokens.
        return lang_1.StringWrapper.replaceAll(content, CR_OR_CRLF_REGEXP, '\n');
    };
    _HtmlTokenizer.prototype.tokenize = function () {
        while (this.peek !== $EOF) {
            var start = this._getLocation();
            try {
                if (this._attemptCharCode($LT)) {
                    if (this._attemptCharCode($BANG)) {
                        if (this._attemptCharCode($LBRACKET)) {
                            this._consumeCdata(start);
                        }
                        else if (this._attemptCharCode($MINUS)) {
                            this._consumeComment(start);
                        }
                        else {
                            this._consumeDocType(start);
                        }
                    }
                    else if (this._attemptCharCode($SLASH)) {
                        this._consumeTagClose(start);
                    }
                    else {
                        this._consumeTagOpen(start);
                    }
                }
                else if (isSpecialFormStart(this.peek, this.nextPeek) && this.tokenizeExpansionForms) {
                    this._consumeExpansionFormStart();
                }
                else if (this.peek === $EQ && this.tokenizeExpansionForms) {
                    this._consumeExpansionCaseStart();
                }
                else if (this.peek === $RBRACE && this.isInExpansionCase() &&
                    this.tokenizeExpansionForms) {
                    this._consumeExpansionCaseEnd();
                }
                else if (this.peek === $RBRACE && this.isInExpansionForm() &&
                    this.tokenizeExpansionForms) {
                    this._consumeExpansionFormEnd();
                }
                else {
                    this._consumeText();
                }
            }
            catch (e) {
                if (e instanceof ControlFlowError) {
                    this.errors.push(e.error);
                }
                else {
                    throw e;
                }
            }
        }
        this._beginToken(HtmlTokenType.EOF);
        this._endToken([]);
        return new HtmlTokenizeResult(mergeTextTokens(this.tokens), this.errors);
    };
    _HtmlTokenizer.prototype._getLocation = function () {
        return new parse_util_1.ParseLocation(this.file, this.index, this.line, this.column);
    };
    _HtmlTokenizer.prototype._getSpan = function (start, end) {
        if (lang_1.isBlank(start)) {
            start = this._getLocation();
        }
        if (lang_1.isBlank(end)) {
            end = this._getLocation();
        }
        return new parse_util_1.ParseSourceSpan(start, end);
    };
    _HtmlTokenizer.prototype._beginToken = function (type, start) {
        if (start === void 0) { start = null; }
        if (lang_1.isBlank(start)) {
            start = this._getLocation();
        }
        this.currentTokenStart = start;
        this.currentTokenType = type;
    };
    _HtmlTokenizer.prototype._endToken = function (parts, end) {
        if (end === void 0) { end = null; }
        if (lang_1.isBlank(end)) {
            end = this._getLocation();
        }
        var token = new HtmlToken(this.currentTokenType, parts, new parse_util_1.ParseSourceSpan(this.currentTokenStart, end));
        this.tokens.push(token);
        this.currentTokenStart = null;
        this.currentTokenType = null;
        return token;
    };
    _HtmlTokenizer.prototype._createError = function (msg, span) {
        var error = new HtmlTokenError(msg, this.currentTokenType, span);
        this.currentTokenStart = null;
        this.currentTokenType = null;
        return new ControlFlowError(error);
    };
    _HtmlTokenizer.prototype._advance = function () {
        if (this.index >= this.length) {
            throw this._createError(unexpectedCharacterErrorMsg($EOF), this._getSpan());
        }
        if (this.peek === $LF) {
            this.line++;
            this.column = 0;
        }
        else if (this.peek !== $LF && this.peek !== $CR) {
            this.column++;
        }
        this.index++;
        this.peek = this.index >= this.length ? $EOF : lang_1.StringWrapper.charCodeAt(this.input, this.index);
        this.nextPeek =
            this.index + 1 >= this.length ? $EOF : lang_1.StringWrapper.charCodeAt(this.input, this.index + 1);
    };
    _HtmlTokenizer.prototype._attemptCharCode = function (charCode) {
        if (this.peek === charCode) {
            this._advance();
            return true;
        }
        return false;
    };
    _HtmlTokenizer.prototype._attemptCharCodeCaseInsensitive = function (charCode) {
        if (compareCharCodeCaseInsensitive(this.peek, charCode)) {
            this._advance();
            return true;
        }
        return false;
    };
    _HtmlTokenizer.prototype._requireCharCode = function (charCode) {
        var location = this._getLocation();
        if (!this._attemptCharCode(charCode)) {
            throw this._createError(unexpectedCharacterErrorMsg(this.peek), this._getSpan(location, location));
        }
    };
    _HtmlTokenizer.prototype._attemptStr = function (chars) {
        for (var i = 0; i < chars.length; i++) {
            if (!this._attemptCharCode(lang_1.StringWrapper.charCodeAt(chars, i))) {
                return false;
            }
        }
        return true;
    };
    _HtmlTokenizer.prototype._attemptStrCaseInsensitive = function (chars) {
        for (var i = 0; i < chars.length; i++) {
            if (!this._attemptCharCodeCaseInsensitive(lang_1.StringWrapper.charCodeAt(chars, i))) {
                return false;
            }
        }
        return true;
    };
    _HtmlTokenizer.prototype._requireStr = function (chars) {
        var location = this._getLocation();
        if (!this._attemptStr(chars)) {
            throw this._createError(unexpectedCharacterErrorMsg(this.peek), this._getSpan(location));
        }
    };
    _HtmlTokenizer.prototype._attemptCharCodeUntilFn = function (predicate) {
        while (!predicate(this.peek)) {
            this._advance();
        }
    };
    _HtmlTokenizer.prototype._requireCharCodeUntilFn = function (predicate, len) {
        var start = this._getLocation();
        this._attemptCharCodeUntilFn(predicate);
        if (this.index - start.offset < len) {
            throw this._createError(unexpectedCharacterErrorMsg(this.peek), this._getSpan(start, start));
        }
    };
    _HtmlTokenizer.prototype._attemptUntilChar = function (char) {
        while (this.peek !== char) {
            this._advance();
        }
    };
    _HtmlTokenizer.prototype._readChar = function (decodeEntities) {
        if (decodeEntities && this.peek === $AMPERSAND) {
            return this._decodeEntity();
        }
        else {
            var index = this.index;
            this._advance();
            return this.input[index];
        }
    };
    _HtmlTokenizer.prototype._decodeEntity = function () {
        var start = this._getLocation();
        this._advance();
        if (this._attemptCharCode($HASH)) {
            var isHex = this._attemptCharCode($x) || this._attemptCharCode($X);
            var numberStart = this._getLocation().offset;
            this._attemptCharCodeUntilFn(isDigitEntityEnd);
            if (this.peek != $SEMICOLON) {
                throw this._createError(unexpectedCharacterErrorMsg(this.peek), this._getSpan());
            }
            this._advance();
            var strNum = this.input.substring(numberStart, this.index - 1);
            try {
                var charCode = lang_1.NumberWrapper.parseInt(strNum, isHex ? 16 : 10);
                return lang_1.StringWrapper.fromCharCode(charCode);
            }
            catch (e) {
                var entity = this.input.substring(start.offset + 1, this.index - 1);
                throw this._createError(unknownEntityErrorMsg(entity), this._getSpan(start));
            }
        }
        else {
            var startPosition = this._savePosition();
            this._attemptCharCodeUntilFn(isNamedEntityEnd);
            if (this.peek != $SEMICOLON) {
                this._restorePosition(startPosition);
                return '&';
            }
            this._advance();
            var name_1 = this.input.substring(start.offset + 1, this.index - 1);
            var char = html_tags_1.NAMED_ENTITIES[name_1];
            if (lang_1.isBlank(char)) {
                throw this._createError(unknownEntityErrorMsg(name_1), this._getSpan(start));
            }
            return char;
        }
    };
    _HtmlTokenizer.prototype._consumeRawText = function (decodeEntities, firstCharOfEnd, attemptEndRest) {
        var tagCloseStart;
        var textStart = this._getLocation();
        this._beginToken(decodeEntities ? HtmlTokenType.ESCAPABLE_RAW_TEXT : HtmlTokenType.RAW_TEXT, textStart);
        var parts = [];
        while (true) {
            tagCloseStart = this._getLocation();
            if (this._attemptCharCode(firstCharOfEnd) && attemptEndRest()) {
                break;
            }
            if (this.index > tagCloseStart.offset) {
                parts.push(this.input.substring(tagCloseStart.offset, this.index));
            }
            while (this.peek !== firstCharOfEnd) {
                parts.push(this._readChar(decodeEntities));
            }
        }
        return this._endToken([this._processCarriageReturns(parts.join(''))], tagCloseStart);
    };
    _HtmlTokenizer.prototype._consumeComment = function (start) {
        var _this = this;
        this._beginToken(HtmlTokenType.COMMENT_START, start);
        this._requireCharCode($MINUS);
        this._endToken([]);
        var textToken = this._consumeRawText(false, $MINUS, function () { return _this._attemptStr('->'); });
        this._beginToken(HtmlTokenType.COMMENT_END, textToken.sourceSpan.end);
        this._endToken([]);
    };
    _HtmlTokenizer.prototype._consumeCdata = function (start) {
        var _this = this;
        this._beginToken(HtmlTokenType.CDATA_START, start);
        this._requireStr('CDATA[');
        this._endToken([]);
        var textToken = this._consumeRawText(false, $RBRACKET, function () { return _this._attemptStr(']>'); });
        this._beginToken(HtmlTokenType.CDATA_END, textToken.sourceSpan.end);
        this._endToken([]);
    };
    _HtmlTokenizer.prototype._consumeDocType = function (start) {
        this._beginToken(HtmlTokenType.DOC_TYPE, start);
        this._attemptUntilChar($GT);
        this._advance();
        this._endToken([this.input.substring(start.offset + 2, this.index - 1)]);
    };
    _HtmlTokenizer.prototype._consumePrefixAndName = function () {
        var nameOrPrefixStart = this.index;
        var prefix = null;
        while (this.peek !== $COLON && !isPrefixEnd(this.peek)) {
            this._advance();
        }
        var nameStart;
        if (this.peek === $COLON) {
            this._advance();
            prefix = this.input.substring(nameOrPrefixStart, this.index - 1);
            nameStart = this.index;
        }
        else {
            nameStart = nameOrPrefixStart;
        }
        this._requireCharCodeUntilFn(isNameEnd, this.index === nameStart ? 1 : 0);
        var name = this.input.substring(nameStart, this.index);
        return [prefix, name];
    };
    _HtmlTokenizer.prototype._consumeTagOpen = function (start) {
        var savedPos = this._savePosition();
        var lowercaseTagName;
        try {
            if (!isAsciiLetter(this.peek)) {
                throw this._createError(unexpectedCharacterErrorMsg(this.peek), this._getSpan());
            }
            var nameStart = this.index;
            this._consumeTagOpenStart(start);
            lowercaseTagName = this.input.substring(nameStart, this.index).toLowerCase();
            this._attemptCharCodeUntilFn(isNotWhitespace);
            while (this.peek !== $SLASH && this.peek !== $GT) {
                this._consumeAttributeName();
                this._attemptCharCodeUntilFn(isNotWhitespace);
                if (this._attemptCharCode($EQ)) {
                    this._attemptCharCodeUntilFn(isNotWhitespace);
                    this._consumeAttributeValue();
                }
                this._attemptCharCodeUntilFn(isNotWhitespace);
            }
            this._consumeTagOpenEnd();
        }
        catch (e) {
            if (e instanceof ControlFlowError) {
                // When the start tag is invalid, assume we want a "<"
                this._restorePosition(savedPos);
                // Back to back text tokens are merged at the end
                this._beginToken(HtmlTokenType.TEXT, start);
                this._endToken(['<']);
                return;
            }
            throw e;
        }
        var contentTokenType = html_tags_1.getHtmlTagDefinition(lowercaseTagName).contentType;
        if (contentTokenType === html_tags_1.HtmlTagContentType.RAW_TEXT) {
            this._consumeRawTextWithTagClose(lowercaseTagName, false);
        }
        else if (contentTokenType === html_tags_1.HtmlTagContentType.ESCAPABLE_RAW_TEXT) {
            this._consumeRawTextWithTagClose(lowercaseTagName, true);
        }
    };
    _HtmlTokenizer.prototype._consumeRawTextWithTagClose = function (lowercaseTagName, decodeEntities) {
        var _this = this;
        var textToken = this._consumeRawText(decodeEntities, $LT, function () {
            if (!_this._attemptCharCode($SLASH))
                return false;
            _this._attemptCharCodeUntilFn(isNotWhitespace);
            if (!_this._attemptStrCaseInsensitive(lowercaseTagName))
                return false;
            _this._attemptCharCodeUntilFn(isNotWhitespace);
            if (!_this._attemptCharCode($GT))
                return false;
            return true;
        });
        this._beginToken(HtmlTokenType.TAG_CLOSE, textToken.sourceSpan.end);
        this._endToken([null, lowercaseTagName]);
    };
    _HtmlTokenizer.prototype._consumeTagOpenStart = function (start) {
        this._beginToken(HtmlTokenType.TAG_OPEN_START, start);
        var parts = this._consumePrefixAndName();
        this._endToken(parts);
    };
    _HtmlTokenizer.prototype._consumeAttributeName = function () {
        this._beginToken(HtmlTokenType.ATTR_NAME);
        var prefixAndName = this._consumePrefixAndName();
        this._endToken(prefixAndName);
    };
    _HtmlTokenizer.prototype._consumeAttributeValue = function () {
        this._beginToken(HtmlTokenType.ATTR_VALUE);
        var value;
        if (this.peek === $SQ || this.peek === $DQ) {
            var quoteChar = this.peek;
            this._advance();
            var parts = [];
            while (this.peek !== quoteChar) {
                parts.push(this._readChar(true));
            }
            value = parts.join('');
            this._advance();
        }
        else {
            var valueStart = this.index;
            this._requireCharCodeUntilFn(isNameEnd, 1);
            value = this.input.substring(valueStart, this.index);
        }
        this._endToken([this._processCarriageReturns(value)]);
    };
    _HtmlTokenizer.prototype._consumeTagOpenEnd = function () {
        var tokenType = this._attemptCharCode($SLASH) ? HtmlTokenType.TAG_OPEN_END_VOID :
            HtmlTokenType.TAG_OPEN_END;
        this._beginToken(tokenType);
        this._requireCharCode($GT);
        this._endToken([]);
    };
    _HtmlTokenizer.prototype._consumeTagClose = function (start) {
        this._beginToken(HtmlTokenType.TAG_CLOSE, start);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        var prefixAndName;
        prefixAndName = this._consumePrefixAndName();
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._requireCharCode($GT);
        this._endToken(prefixAndName);
    };
    _HtmlTokenizer.prototype._consumeExpansionFormStart = function () {
        this._beginToken(HtmlTokenType.EXPANSION_FORM_START, this._getLocation());
        this._requireCharCode($LBRACE);
        this._endToken([]);
        this._beginToken(HtmlTokenType.RAW_TEXT, this._getLocation());
        var condition = this._readUntil($COMMA);
        this._endToken([condition], this._getLocation());
        this._requireCharCode($COMMA);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._beginToken(HtmlTokenType.RAW_TEXT, this._getLocation());
        var type = this._readUntil($COMMA);
        this._endToken([type], this._getLocation());
        this._requireCharCode($COMMA);
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this.expansionCaseStack.push(HtmlTokenType.EXPANSION_FORM_START);
    };
    _HtmlTokenizer.prototype._consumeExpansionCaseStart = function () {
        this._requireCharCode($EQ);
        this._beginToken(HtmlTokenType.EXPANSION_CASE_VALUE, this._getLocation());
        var value = this._readUntil($LBRACE).trim();
        this._endToken([value], this._getLocation());
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this._beginToken(HtmlTokenType.EXPANSION_CASE_EXP_START, this._getLocation());
        this._requireCharCode($LBRACE);
        this._endToken([], this._getLocation());
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this.expansionCaseStack.push(HtmlTokenType.EXPANSION_CASE_EXP_START);
    };
    _HtmlTokenizer.prototype._consumeExpansionCaseEnd = function () {
        this._beginToken(HtmlTokenType.EXPANSION_CASE_EXP_END, this._getLocation());
        this._requireCharCode($RBRACE);
        this._endToken([], this._getLocation());
        this._attemptCharCodeUntilFn(isNotWhitespace);
        this.expansionCaseStack.pop();
    };
    _HtmlTokenizer.prototype._consumeExpansionFormEnd = function () {
        this._beginToken(HtmlTokenType.EXPANSION_FORM_END, this._getLocation());
        this._requireCharCode($RBRACE);
        this._endToken([]);
        this.expansionCaseStack.pop();
    };
    _HtmlTokenizer.prototype._consumeText = function () {
        var start = this._getLocation();
        this._beginToken(HtmlTokenType.TEXT, start);
        var parts = [];
        var interpolation = false;
        if (this.peek === $LBRACE && this.nextPeek === $LBRACE) {
            parts.push(this._readChar(true));
            parts.push(this._readChar(true));
            interpolation = true;
        }
        else {
            parts.push(this._readChar(true));
        }
        while (!this.isTextEnd(interpolation)) {
            if (this.peek === $LBRACE && this.nextPeek === $LBRACE) {
                parts.push(this._readChar(true));
                parts.push(this._readChar(true));
                interpolation = true;
            }
            else if (this.peek === $RBRACE && this.nextPeek === $RBRACE && interpolation) {
                parts.push(this._readChar(true));
                parts.push(this._readChar(true));
                interpolation = false;
            }
            else {
                parts.push(this._readChar(true));
            }
        }
        this._endToken([this._processCarriageReturns(parts.join(''))]);
    };
    _HtmlTokenizer.prototype.isTextEnd = function (interpolation) {
        if (this.peek === $LT || this.peek === $EOF)
            return true;
        if (this.tokenizeExpansionForms) {
            if (isSpecialFormStart(this.peek, this.nextPeek))
                return true;
            if (this.peek === $RBRACE && !interpolation &&
                (this.isInExpansionCase() || this.isInExpansionForm()))
                return true;
        }
        return false;
    };
    _HtmlTokenizer.prototype._savePosition = function () {
        return [this.peek, this.index, this.column, this.line, this.tokens.length];
    };
    _HtmlTokenizer.prototype._readUntil = function (char) {
        var start = this.index;
        this._attemptUntilChar(char);
        return this.input.substring(start, this.index);
    };
    _HtmlTokenizer.prototype._restorePosition = function (position) {
        this.peek = position[0];
        this.index = position[1];
        this.column = position[2];
        this.line = position[3];
        var nbTokens = position[4];
        if (nbTokens < this.tokens.length) {
            // remove any extra tokens
            this.tokens = collection_1.ListWrapper.slice(this.tokens, 0, nbTokens);
        }
    };
    _HtmlTokenizer.prototype.isInExpansionCase = function () {
        return this.expansionCaseStack.length > 0 &&
            this.expansionCaseStack[this.expansionCaseStack.length - 1] ===
                HtmlTokenType.EXPANSION_CASE_EXP_START;
    };
    _HtmlTokenizer.prototype.isInExpansionForm = function () {
        return this.expansionCaseStack.length > 0 &&
            this.expansionCaseStack[this.expansionCaseStack.length - 1] ===
                HtmlTokenType.EXPANSION_FORM_START;
    };
    return _HtmlTokenizer;
}());
function isNotWhitespace(code) {
    return !isWhitespace(code) || code === $EOF;
}
function isWhitespace(code) {
    return (code >= $TAB && code <= $SPACE) || (code === $NBSP);
}
function isNameEnd(code) {
    return isWhitespace(code) || code === $GT || code === $SLASH || code === $SQ || code === $DQ ||
        code === $EQ;
}
function isPrefixEnd(code) {
    return (code < $a || $z < code) && (code < $A || $Z < code) && (code < $0 || code > $9);
}
function isDigitEntityEnd(code) {
    return code == $SEMICOLON || code == $EOF || !isAsciiHexDigit(code);
}
function isNamedEntityEnd(code) {
    return code == $SEMICOLON || code == $EOF || !isAsciiLetter(code);
}
function isSpecialFormStart(peek, nextPeek) {
    return peek === $LBRACE && nextPeek != $LBRACE;
}
function isAsciiLetter(code) {
    return code >= $a && code <= $z || code >= $A && code <= $Z;
}
function isAsciiHexDigit(code) {
    return code >= $a && code <= $f || code >= $A && code <= $F || code >= $0 && code <= $9;
}
function compareCharCodeCaseInsensitive(code1, code2) {
    return toUpperCaseCharCode(code1) == toUpperCaseCharCode(code2);
}
function toUpperCaseCharCode(code) {
    return code >= $a && code <= $z ? code - $a + $A : code;
}
function mergeTextTokens(srcTokens) {
    var dstTokens = [];
    var lastDstToken;
    for (var i = 0; i < srcTokens.length; i++) {
        var token = srcTokens[i];
        if (lang_1.isPresent(lastDstToken) && lastDstToken.type == HtmlTokenType.TEXT &&
            token.type == HtmlTokenType.TEXT) {
            lastDstToken.parts[0] += token.parts[0];
            lastDstToken.sourceSpan.end = token.sourceSpan.end;
        }
        else {
            lastDstToken = token;
            dstTokens.push(lastDstToken);
        }
    }
    return dstTokens;
}
//# sourceMappingURL=html_lexer.js.map