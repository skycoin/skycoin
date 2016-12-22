"use strict";
var lang_1 = require('../src/facade/lang');
var html_tags_1 = require('./html_tags');
var NG_CONTENT_SELECT_ATTR = 'select';
var NG_CONTENT_ELEMENT = 'ng-content';
var LINK_ELEMENT = 'link';
var LINK_STYLE_REL_ATTR = 'rel';
var LINK_STYLE_HREF_ATTR = 'href';
var LINK_STYLE_REL_VALUE = 'stylesheet';
var STYLE_ELEMENT = 'style';
var SCRIPT_ELEMENT = 'script';
var NG_NON_BINDABLE_ATTR = 'ngNonBindable';
var NG_PROJECT_AS = 'ngProjectAs';
function preparseElement(ast) {
    var selectAttr = null;
    var hrefAttr = null;
    var relAttr = null;
    var nonBindable = false;
    var projectAs = null;
    ast.attrs.forEach(function (attr) {
        var lcAttrName = attr.name.toLowerCase();
        if (lcAttrName == NG_CONTENT_SELECT_ATTR) {
            selectAttr = attr.value;
        }
        else if (lcAttrName == LINK_STYLE_HREF_ATTR) {
            hrefAttr = attr.value;
        }
        else if (lcAttrName == LINK_STYLE_REL_ATTR) {
            relAttr = attr.value;
        }
        else if (attr.name == NG_NON_BINDABLE_ATTR) {
            nonBindable = true;
        }
        else if (attr.name == NG_PROJECT_AS) {
            if (attr.value.length > 0) {
                projectAs = attr.value;
            }
        }
    });
    selectAttr = normalizeNgContentSelect(selectAttr);
    var nodeName = ast.name.toLowerCase();
    var type = PreparsedElementType.OTHER;
    if (html_tags_1.splitNsName(nodeName)[1] == NG_CONTENT_ELEMENT) {
        type = PreparsedElementType.NG_CONTENT;
    }
    else if (nodeName == STYLE_ELEMENT) {
        type = PreparsedElementType.STYLE;
    }
    else if (nodeName == SCRIPT_ELEMENT) {
        type = PreparsedElementType.SCRIPT;
    }
    else if (nodeName == LINK_ELEMENT && relAttr == LINK_STYLE_REL_VALUE) {
        type = PreparsedElementType.STYLESHEET;
    }
    return new PreparsedElement(type, selectAttr, hrefAttr, nonBindable, projectAs);
}
exports.preparseElement = preparseElement;
(function (PreparsedElementType) {
    PreparsedElementType[PreparsedElementType["NG_CONTENT"] = 0] = "NG_CONTENT";
    PreparsedElementType[PreparsedElementType["STYLE"] = 1] = "STYLE";
    PreparsedElementType[PreparsedElementType["STYLESHEET"] = 2] = "STYLESHEET";
    PreparsedElementType[PreparsedElementType["SCRIPT"] = 3] = "SCRIPT";
    PreparsedElementType[PreparsedElementType["OTHER"] = 4] = "OTHER";
})(exports.PreparsedElementType || (exports.PreparsedElementType = {}));
var PreparsedElementType = exports.PreparsedElementType;
var PreparsedElement = (function () {
    function PreparsedElement(type, selectAttr, hrefAttr, nonBindable, projectAs) {
        this.type = type;
        this.selectAttr = selectAttr;
        this.hrefAttr = hrefAttr;
        this.nonBindable = nonBindable;
        this.projectAs = projectAs;
    }
    return PreparsedElement;
}());
exports.PreparsedElement = PreparsedElement;
function normalizeNgContentSelect(selectAttr) {
    if (lang_1.isBlank(selectAttr) || selectAttr.length === 0) {
        return '*';
    }
    return selectAttr;
}
//# sourceMappingURL=template_preparser.js.map