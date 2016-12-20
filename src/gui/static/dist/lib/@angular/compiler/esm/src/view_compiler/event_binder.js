import { isBlank, isPresent, StringWrapper } from '../../src/facade/lang';
import { ListWrapper, StringMapWrapper } from '../../src/facade/collection';
import { EventHandlerVars, ViewProperties } from './constants';
import * as o from '../output/output_ast';
import { CompileMethod } from './compile_method';
import { convertCdStatementToIr } from './expression_converter';
import { CompileBinding } from './compile_binding';
export class CompileEventListener {
    constructor(compileElement, eventTarget, eventName, listenerIndex) {
        this.compileElement = compileElement;
        this.eventTarget = eventTarget;
        this.eventName = eventName;
        this._hasComponentHostListener = false;
        this._actionResultExprs = [];
        this._method = new CompileMethod(compileElement.view);
        this._methodName =
            `_handle_${santitizeEventName(eventName)}_${compileElement.nodeIndex}_${listenerIndex}`;
        this._eventParam =
            new o.FnParam(EventHandlerVars.event.name, o.importType(this.compileElement.view.genConfig.renderTypes.renderEvent));
    }
    static getOrCreate(compileElement, eventTarget, eventName, targetEventListeners) {
        var listener = targetEventListeners.find(listener => listener.eventTarget == eventTarget &&
            listener.eventName == eventName);
        if (isBlank(listener)) {
            listener = new CompileEventListener(compileElement, eventTarget, eventName, targetEventListeners.length);
            targetEventListeners.push(listener);
        }
        return listener;
    }
    addAction(hostEvent, directive, directiveInstance) {
        if (isPresent(directive) && directive.isComponent) {
            this._hasComponentHostListener = true;
        }
        this._method.resetDebugInfo(this.compileElement.nodeIndex, hostEvent);
        var context = isPresent(directiveInstance) ? directiveInstance :
            this.compileElement.view.componentContext;
        var actionStmts = convertCdStatementToIr(this.compileElement.view, context, hostEvent.handler);
        var lastIndex = actionStmts.length - 1;
        if (lastIndex >= 0) {
            var lastStatement = actionStmts[lastIndex];
            var returnExpr = convertStmtIntoExpression(lastStatement);
            var preventDefaultVar = o.variable(`pd_${this._actionResultExprs.length}`);
            this._actionResultExprs.push(preventDefaultVar);
            if (isPresent(returnExpr)) {
                // Note: We need to cast the result of the method call to dynamic,
                // as it might be a void method!
                actionStmts[lastIndex] =
                    preventDefaultVar.set(returnExpr.cast(o.DYNAMIC_TYPE).notIdentical(o.literal(false)))
                        .toDeclStmt(null, [o.StmtModifier.Final]);
            }
        }
        this._method.addStmts(actionStmts);
    }
    finishMethod() {
        var markPathToRootStart = this._hasComponentHostListener ?
            this.compileElement.appElement.prop('componentView') :
            o.THIS_EXPR;
        var resultExpr = o.literal(true);
        this._actionResultExprs.forEach((expr) => { resultExpr = resultExpr.and(expr); });
        var stmts = [markPathToRootStart.callMethod('markPathToRootAsCheckOnce', []).toStmt()]
            .concat(this._method.finish())
            .concat([new o.ReturnStatement(resultExpr)]);
        // private is fine here as no child view will reference the event handler...
        this.compileElement.view.eventHandlerMethods.push(new o.ClassMethod(this._methodName, [this._eventParam], stmts, o.BOOL_TYPE, [o.StmtModifier.Private]));
    }
    listenToRenderer() {
        var listenExpr;
        var eventListener = o.THIS_EXPR.callMethod('eventHandler', [o.THIS_EXPR.prop(this._methodName).callMethod(o.BuiltinMethod.bind, [o.THIS_EXPR])]);
        if (isPresent(this.eventTarget)) {
            listenExpr = ViewProperties.renderer.callMethod('listenGlobal', [o.literal(this.eventTarget), o.literal(this.eventName), eventListener]);
        }
        else {
            listenExpr = ViewProperties.renderer.callMethod('listen', [this.compileElement.renderNode, o.literal(this.eventName), eventListener]);
        }
        var disposable = o.variable(`disposable_${this.compileElement.view.disposables.length}`);
        this.compileElement.view.disposables.push(disposable);
        // private is fine here as no child view will reference the event handler...
        this.compileElement.view.createMethod.addStmt(disposable.set(listenExpr).toDeclStmt(o.FUNCTION_TYPE, [o.StmtModifier.Private]));
    }
    listenToDirective(directiveInstance, observablePropName) {
        var subscription = o.variable(`subscription_${this.compileElement.view.subscriptions.length}`);
        this.compileElement.view.subscriptions.push(subscription);
        var eventListener = o.THIS_EXPR.callMethod('eventHandler', [o.THIS_EXPR.prop(this._methodName).callMethod(o.BuiltinMethod.bind, [o.THIS_EXPR])]);
        this.compileElement.view.createMethod.addStmt(subscription.set(directiveInstance.prop(observablePropName)
            .callMethod(o.BuiltinMethod.SubscribeObservable, [eventListener]))
            .toDeclStmt(null, [o.StmtModifier.Final]));
    }
}
export function collectEventListeners(hostEvents, dirs, compileElement) {
    var eventListeners = [];
    hostEvents.forEach((hostEvent) => {
        compileElement.view.bindings.push(new CompileBinding(compileElement, hostEvent));
        var listener = CompileEventListener.getOrCreate(compileElement, hostEvent.target, hostEvent.name, eventListeners);
        listener.addAction(hostEvent, null, null);
    });
    ListWrapper.forEachWithIndex(dirs, (directiveAst, i) => {
        var directiveInstance = compileElement.directiveInstances[i];
        directiveAst.hostEvents.forEach((hostEvent) => {
            compileElement.view.bindings.push(new CompileBinding(compileElement, hostEvent));
            var listener = CompileEventListener.getOrCreate(compileElement, hostEvent.target, hostEvent.name, eventListeners);
            listener.addAction(hostEvent, directiveAst.directive, directiveInstance);
        });
    });
    eventListeners.forEach((listener) => listener.finishMethod());
    return eventListeners;
}
export function bindDirectiveOutputs(directiveAst, directiveInstance, eventListeners) {
    StringMapWrapper.forEach(directiveAst.directive.outputs, (eventName, observablePropName) => {
        eventListeners.filter(listener => listener.eventName == eventName)
            .forEach((listener) => { listener.listenToDirective(directiveInstance, observablePropName); });
    });
}
export function bindRenderOutputs(eventListeners) {
    eventListeners.forEach(listener => listener.listenToRenderer());
}
function convertStmtIntoExpression(stmt) {
    if (stmt instanceof o.ExpressionStatement) {
        return stmt.expr;
    }
    else if (stmt instanceof o.ReturnStatement) {
        return stmt.value;
    }
    return null;
}
function santitizeEventName(name) {
    return StringWrapper.replaceAll(name, /[^a-zA-Z_]/g, '_');
}
//# sourceMappingURL=event_binder.js.map