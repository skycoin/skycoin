import { Inject, Injectable, ViewEncapsulation } from '@angular/core';
import { AnimationBuilder } from '../animate/animation_builder';
import { isPresent, isBlank, Json, RegExpWrapper, stringify, StringWrapper, isArray, isString } from '../../src/facade/lang';
import { BaseException } from '../../src/facade/exceptions';
import { DomSharedStylesHost } from './shared_styles_host';
import { EventManager } from './events/event_manager';
import { DOCUMENT } from './dom_tokens';
import { getDOM } from './dom_adapter';
import { camelCaseToDashCase } from './util';
const NAMESPACE_URIS = 
/*@ts2dart_const*/
{ 'xlink': 'http://www.w3.org/1999/xlink', 'svg': 'http://www.w3.org/2000/svg' };
const TEMPLATE_COMMENT_TEXT = 'template bindings={}';
var TEMPLATE_BINDINGS_EXP = /^template bindings=(.*)$/g;
export class DomRootRenderer {
    constructor(document, eventManager, sharedStylesHost, animate) {
        this.document = document;
        this.eventManager = eventManager;
        this.sharedStylesHost = sharedStylesHost;
        this.animate = animate;
        this._registeredComponents = new Map();
    }
    renderComponent(componentProto) {
        var renderer = this._registeredComponents.get(componentProto.id);
        if (isBlank(renderer)) {
            renderer = new DomRenderer(this, componentProto);
            this._registeredComponents.set(componentProto.id, renderer);
        }
        return renderer;
    }
}
export class DomRootRenderer_ extends DomRootRenderer {
    constructor(_document, _eventManager, sharedStylesHost, animate) {
        super(_document, _eventManager, sharedStylesHost, animate);
    }
}
DomRootRenderer_.decorators = [
    { type: Injectable },
];
DomRootRenderer_.ctorParameters = [
    { type: undefined, decorators: [{ type: Inject, args: [DOCUMENT,] },] },
    { type: EventManager, },
    { type: DomSharedStylesHost, },
    { type: AnimationBuilder, },
];
export class DomRenderer {
    constructor(_rootRenderer, componentProto) {
        this._rootRenderer = _rootRenderer;
        this.componentProto = componentProto;
        this._styles = _flattenStyles(componentProto.id, componentProto.styles, []);
        if (componentProto.encapsulation !== ViewEncapsulation.Native) {
            this._rootRenderer.sharedStylesHost.addStyles(this._styles);
        }
        if (this.componentProto.encapsulation === ViewEncapsulation.Emulated) {
            this._contentAttr = _shimContentAttribute(componentProto.id);
            this._hostAttr = _shimHostAttribute(componentProto.id);
        }
        else {
            this._contentAttr = null;
            this._hostAttr = null;
        }
    }
    selectRootElement(selectorOrNode, debugInfo) {
        var el;
        if (isString(selectorOrNode)) {
            el = getDOM().querySelector(this._rootRenderer.document, selectorOrNode);
            if (isBlank(el)) {
                throw new BaseException(`The selector "${selectorOrNode}" did not match any elements`);
            }
        }
        else {
            el = selectorOrNode;
        }
        getDOM().clearNodes(el);
        return el;
    }
    createElement(parent, name, debugInfo) {
        var nsAndName = splitNamespace(name);
        var el = isPresent(nsAndName[0]) ?
            getDOM().createElementNS(NAMESPACE_URIS[nsAndName[0]], nsAndName[1]) :
            getDOM().createElement(nsAndName[1]);
        if (isPresent(this._contentAttr)) {
            getDOM().setAttribute(el, this._contentAttr, '');
        }
        if (isPresent(parent)) {
            getDOM().appendChild(parent, el);
        }
        return el;
    }
    createViewRoot(hostElement) {
        var nodesParent;
        if (this.componentProto.encapsulation === ViewEncapsulation.Native) {
            nodesParent = getDOM().createShadowRoot(hostElement);
            this._rootRenderer.sharedStylesHost.addHost(nodesParent);
            for (var i = 0; i < this._styles.length; i++) {
                getDOM().appendChild(nodesParent, getDOM().createStyleElement(this._styles[i]));
            }
        }
        else {
            if (isPresent(this._hostAttr)) {
                getDOM().setAttribute(hostElement, this._hostAttr, '');
            }
            nodesParent = hostElement;
        }
        return nodesParent;
    }
    createTemplateAnchor(parentElement, debugInfo) {
        var comment = getDOM().createComment(TEMPLATE_COMMENT_TEXT);
        if (isPresent(parentElement)) {
            getDOM().appendChild(parentElement, comment);
        }
        return comment;
    }
    createText(parentElement, value, debugInfo) {
        var node = getDOM().createTextNode(value);
        if (isPresent(parentElement)) {
            getDOM().appendChild(parentElement, node);
        }
        return node;
    }
    projectNodes(parentElement, nodes) {
        if (isBlank(parentElement))
            return;
        appendNodes(parentElement, nodes);
    }
    attachViewAfter(node, viewRootNodes) {
        moveNodesAfterSibling(node, viewRootNodes);
        for (let i = 0; i < viewRootNodes.length; i++)
            this.animateNodeEnter(viewRootNodes[i]);
    }
    detachView(viewRootNodes) {
        for (var i = 0; i < viewRootNodes.length; i++) {
            var node = viewRootNodes[i];
            getDOM().remove(node);
            this.animateNodeLeave(node);
        }
    }
    destroyView(hostElement, viewAllNodes) {
        if (this.componentProto.encapsulation === ViewEncapsulation.Native && isPresent(hostElement)) {
            this._rootRenderer.sharedStylesHost.removeHost(getDOM().getShadowRoot(hostElement));
        }
    }
    listen(renderElement, name, callback) {
        return this._rootRenderer.eventManager.addEventListener(renderElement, name, decoratePreventDefault(callback));
    }
    listenGlobal(target, name, callback) {
        return this._rootRenderer.eventManager.addGlobalEventListener(target, name, decoratePreventDefault(callback));
    }
    setElementProperty(renderElement, propertyName, propertyValue) {
        getDOM().setProperty(renderElement, propertyName, propertyValue);
    }
    setElementAttribute(renderElement, attributeName, attributeValue) {
        var attrNs;
        var nsAndName = splitNamespace(attributeName);
        if (isPresent(nsAndName[0])) {
            attributeName = nsAndName[0] + ':' + nsAndName[1];
            attrNs = NAMESPACE_URIS[nsAndName[0]];
        }
        if (isPresent(attributeValue)) {
            if (isPresent(attrNs)) {
                getDOM().setAttributeNS(renderElement, attrNs, attributeName, attributeValue);
            }
            else {
                getDOM().setAttribute(renderElement, attributeName, attributeValue);
            }
        }
        else {
            if (isPresent(attrNs)) {
                getDOM().removeAttributeNS(renderElement, attrNs, nsAndName[1]);
            }
            else {
                getDOM().removeAttribute(renderElement, attributeName);
            }
        }
    }
    setBindingDebugInfo(renderElement, propertyName, propertyValue) {
        var dashCasedPropertyName = camelCaseToDashCase(propertyName);
        if (getDOM().isCommentNode(renderElement)) {
            var existingBindings = RegExpWrapper.firstMatch(TEMPLATE_BINDINGS_EXP, StringWrapper.replaceAll(getDOM().getText(renderElement), /\n/g, ''));
            var parsedBindings = Json.parse(existingBindings[1]);
            parsedBindings[dashCasedPropertyName] = propertyValue;
            getDOM().setText(renderElement, StringWrapper.replace(TEMPLATE_COMMENT_TEXT, '{}', Json.stringify(parsedBindings)));
        }
        else {
            this.setElementAttribute(renderElement, propertyName, propertyValue);
        }
    }
    setElementClass(renderElement, className, isAdd) {
        if (isAdd) {
            getDOM().addClass(renderElement, className);
        }
        else {
            getDOM().removeClass(renderElement, className);
        }
    }
    setElementStyle(renderElement, styleName, styleValue) {
        if (isPresent(styleValue)) {
            getDOM().setStyle(renderElement, styleName, stringify(styleValue));
        }
        else {
            getDOM().removeStyle(renderElement, styleName);
        }
    }
    invokeElementMethod(renderElement, methodName, args) {
        getDOM().invoke(renderElement, methodName, args);
    }
    setText(renderNode, text) { getDOM().setText(renderNode, text); }
    /**
     * Performs animations if necessary
     * @param node
     */
    animateNodeEnter(node) {
        if (getDOM().isElementNode(node) && getDOM().hasClass(node, 'ng-animate')) {
            getDOM().addClass(node, 'ng-enter');
            this._rootRenderer.animate.css()
                .addAnimationClass('ng-enter-active')
                .start(node)
                .onComplete(() => { getDOM().removeClass(node, 'ng-enter'); });
        }
    }
    /**
     * If animations are necessary, performs animations then removes the element; otherwise, it just
     * removes the element.
     * @param node
     */
    animateNodeLeave(node) {
        if (getDOM().isElementNode(node) && getDOM().hasClass(node, 'ng-animate')) {
            getDOM().addClass(node, 'ng-leave');
            this._rootRenderer.animate.css()
                .addAnimationClass('ng-leave-active')
                .start(node)
                .onComplete(() => {
                getDOM().removeClass(node, 'ng-leave');
                getDOM().remove(node);
            });
        }
        else {
            getDOM().remove(node);
        }
    }
}
function moveNodesAfterSibling(sibling, nodes) {
    var parent = getDOM().parentElement(sibling);
    if (nodes.length > 0 && isPresent(parent)) {
        var nextSibling = getDOM().nextSibling(sibling);
        if (isPresent(nextSibling)) {
            for (var i = 0; i < nodes.length; i++) {
                getDOM().insertBefore(nextSibling, nodes[i]);
            }
        }
        else {
            for (var i = 0; i < nodes.length; i++) {
                getDOM().appendChild(parent, nodes[i]);
            }
        }
    }
}
function appendNodes(parent, nodes) {
    for (var i = 0; i < nodes.length; i++) {
        getDOM().appendChild(parent, nodes[i]);
    }
}
function decoratePreventDefault(eventHandler) {
    return (event) => {
        var allowDefaultBehavior = eventHandler(event);
        if (allowDefaultBehavior === false) {
            // TODO(tbosch): move preventDefault into event plugins...
            getDOM().preventDefault(event);
        }
    };
}
var COMPONENT_REGEX = /%COMP%/g;
export const COMPONENT_VARIABLE = '%COMP%';
export const HOST_ATTR = `_nghost-${COMPONENT_VARIABLE}`;
export const CONTENT_ATTR = `_ngcontent-${COMPONENT_VARIABLE}`;
function _shimContentAttribute(componentShortId) {
    return StringWrapper.replaceAll(CONTENT_ATTR, COMPONENT_REGEX, componentShortId);
}
function _shimHostAttribute(componentShortId) {
    return StringWrapper.replaceAll(HOST_ATTR, COMPONENT_REGEX, componentShortId);
}
function _flattenStyles(compId, styles, target) {
    for (var i = 0; i < styles.length; i++) {
        var style = styles[i];
        if (isArray(style)) {
            _flattenStyles(compId, style, target);
        }
        else {
            style = StringWrapper.replaceAll(style, COMPONENT_REGEX, compId);
            target.push(style);
        }
    }
    return target;
}
var NS_PREFIX_RE = /^@([^:]+):(.+)/g;
function splitNamespace(name) {
    if (name[0] != '@') {
        return [null, name];
    }
    let match = RegExpWrapper.firstMatch(NS_PREFIX_RE, name);
    return [match[1], match[2]];
}
//# sourceMappingURL=dom_renderer.js.map