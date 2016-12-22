import { provide, ReflectiveInjector } from '@angular/core';
import { isBlank, isPresent } from './facade/lang';
import { ListWrapper } from './facade/collection';
import { EventEmitter, PromiseWrapper, ObservableWrapper } from './facade/async';
import { StringMapWrapper } from './facade/collection';
import { BaseException } from '@angular/core';
import { recognize } from './recognize';
import { link } from './link';
import { equalSegments, routeSegmentComponentFactory, RouteSegment, RouteTree, rootNode, TreeNode, UrlSegment } from './segments';
import { hasLifecycleHook } from './lifecycle_reflector';
import { DEFAULT_OUTLET_NAME } from './constants';
/**
 * @internal
 */
export class RouterOutletMap {
    constructor() {
        /** @internal */
        this._outlets = {};
    }
    registerOutlet(name, outlet) { this._outlets[name] = outlet; }
}
/**
 * The `Router` is responsible for mapping URLs to components.
 *
 * You can see the state of the router by inspecting the read-only fields `router.urlTree`
 * and `router.routeTree`.
 */
export class Router {
    /**
     * @internal
     */
    constructor(_rootComponent, _rootComponentType, _componentResolver, _urlSerializer, _routerOutletMap, _location) {
        this._rootComponent = _rootComponent;
        this._rootComponentType = _rootComponentType;
        this._componentResolver = _componentResolver;
        this._urlSerializer = _urlSerializer;
        this._routerOutletMap = _routerOutletMap;
        this._location = _location;
        this._changes = new EventEmitter();
        this._prevTree = this._createInitialTree();
        this._setUpLocationChangeListener();
        this.navigateByUrl(this._location.path());
    }
    /**
     * Returns the current url tree.
     */
    get urlTree() { return this._urlTree; }
    /**
     * Returns the current route tree.
     */
    get routeTree() { return this._prevTree; }
    /**
     * An observable or url changes from the router.
     */
    get changes() { return this._changes; }
    /**
     * Navigate based on the provided url. This navigation is always absolute.
     *
     * ### Usage
     *
     * ```
     * router.navigateByUrl("/team/33/user/11");
     * ```
     */
    navigateByUrl(url) {
        return this._navigate(this._urlSerializer.parse(url));
    }
    /**
     * Navigate based on the provided array of commands and a starting point.
     * If no segment is provided, the navigation is absolute.
     *
     * ### Usage
     *
     * ```
     * router.navigate(['team', 33, 'team', '11], segment);
     * ```
     */
    navigate(commands, segment) {
        return this._navigate(this.createUrlTree(commands, segment));
    }
    /**
     * @internal
     */
    dispose() { ObservableWrapper.dispose(this._locationSubscription); }
    /**
     * Applies an array of commands to the current url tree and creates
     * a new url tree.
     *
     * When given a segment, applies the given commands starting from the segment.
     * When not given a segment, applies the given command starting from the root.
     *
     * ### Usage
     *
     * ```
     * // create /team/33/user/11
     * router.createUrlTree(['/team', 33, 'user', 11]);
     *
     * // create /team/33;expand=true/user/11
     * router.createUrlTree(['/team', 33, {expand: true}, 'user', 11]);
     *
     * // you can collapse static fragments like this
     * router.createUrlTree(['/team/33/user', userId]);
     *
     * // assuming the current url is `/team/33/user/11` and the segment points to `user/11`
     *
     * // navigate to /team/33/user/11/details
     * router.createUrlTree(['details'], segment);
     *
     * // navigate to /team/33/user/22
     * router.createUrlTree(['../22'], segment);
     *
     * // navigate to /team/44/user/22
     * router.createUrlTree(['../../team/44/user/22'], segment);
     * ```
     */
    createUrlTree(commands, segment) {
        let s = isPresent(segment) ? segment : this._prevTree.root;
        return link(s, this._prevTree, this.urlTree, commands);
    }
    /**
     * Serializes a {@link UrlTree} into a string.
     */
    serializeUrl(url) { return this._urlSerializer.serialize(url); }
    _createInitialTree() {
        let root = new RouteSegment([new UrlSegment("", {}, null)], {}, DEFAULT_OUTLET_NAME, this._rootComponentType, null);
        return new RouteTree(new TreeNode(root, []));
    }
    _setUpLocationChangeListener() {
        this._locationSubscription = this._location.subscribe((change) => { this._navigate(this._urlSerializer.parse(change['url'])); });
    }
    _navigate(url) {
        this._urlTree = url;
        return recognize(this._componentResolver, this._rootComponentType, url)
            .then(currTree => {
            return new _LoadSegments(currTree, this._prevTree)
                .load(this._routerOutletMap, this._rootComponent)
                .then(updated => {
                if (updated) {
                    this._prevTree = currTree;
                    this._location.go(this._urlSerializer.serialize(this._urlTree));
                    this._changes.emit(null);
                }
            });
        });
    }
}
class _LoadSegments {
    constructor(currTree, prevTree) {
        this.currTree = currTree;
        this.prevTree = prevTree;
        this.deactivations = [];
        this.performMutation = true;
    }
    load(parentOutletMap, rootComponent) {
        let prevRoot = isPresent(this.prevTree) ? rootNode(this.prevTree) : null;
        let currRoot = rootNode(this.currTree);
        return this.canDeactivate(currRoot, prevRoot, parentOutletMap, rootComponent)
            .then(res => {
            this.performMutation = true;
            if (res) {
                this.loadChildSegments(currRoot, prevRoot, parentOutletMap, [rootComponent]);
            }
            return res;
        });
    }
    canDeactivate(currRoot, prevRoot, outletMap, rootComponent) {
        this.performMutation = false;
        this.loadChildSegments(currRoot, prevRoot, outletMap, [rootComponent]);
        let allPaths = PromiseWrapper.all(this.deactivations.map(r => this.checkCanDeactivatePath(r)));
        return allPaths.then((values) => values.filter(v => v).length === values.length);
    }
    checkCanDeactivatePath(path) {
        let curr = PromiseWrapper.resolve(true);
        for (let p of ListWrapper.reversed(path)) {
            curr = curr.then(_ => {
                if (hasLifecycleHook("routerCanDeactivate", p)) {
                    return p.routerCanDeactivate(this.prevTree, this.currTree);
                }
                else {
                    return _;
                }
            });
        }
        return curr;
    }
    loadChildSegments(currNode, prevNode, outletMap, components) {
        let prevChildren = isPresent(prevNode) ?
            prevNode.children.reduce((m, c) => {
                m[c.value.outlet] = c;
                return m;
            }, {}) :
            {};
        currNode.children.forEach(c => {
            this.loadSegments(c, prevChildren[c.value.outlet], outletMap, components);
            StringMapWrapper.delete(prevChildren, c.value.outlet);
        });
        StringMapWrapper.forEach(prevChildren, (v, k) => this.unloadOutlet(outletMap._outlets[k], components));
    }
    loadSegments(currNode, prevNode, parentOutletMap, components) {
        let curr = currNode.value;
        let prev = isPresent(prevNode) ? prevNode.value : null;
        let outlet = this.getOutlet(parentOutletMap, currNode.value);
        if (equalSegments(curr, prev)) {
            this.loadChildSegments(currNode, prevNode, outlet.outletMap, components.concat([outlet.loadedComponent]));
        }
        else {
            this.unloadOutlet(outlet, components);
            if (this.performMutation) {
                let outletMap = new RouterOutletMap();
                let loadedComponent = this.loadNewSegment(outletMap, curr, prev, outlet);
                this.loadChildSegments(currNode, prevNode, outletMap, components.concat([loadedComponent]));
            }
        }
    }
    loadNewSegment(outletMap, curr, prev, outlet) {
        let resolved = ReflectiveInjector.resolve([provide(RouterOutletMap, { useValue: outletMap }), provide(RouteSegment, { useValue: curr })]);
        let ref = outlet.load(routeSegmentComponentFactory(curr), resolved, outletMap);
        if (hasLifecycleHook("routerOnActivate", ref.instance)) {
            ref.instance.routerOnActivate(curr, prev, this.currTree, this.prevTree);
        }
        return ref.instance;
    }
    getOutlet(outletMap, segment) {
        let outlet = outletMap._outlets[segment.outlet];
        if (isBlank(outlet)) {
            if (segment.outlet == DEFAULT_OUTLET_NAME) {
                throw new BaseException(`Cannot find default outlet`);
            }
            else {
                throw new BaseException(`Cannot find the outlet ${segment.outlet}`);
            }
        }
        return outlet;
    }
    unloadOutlet(outlet, components) {
        if (isPresent(outlet) && outlet.isLoaded) {
            StringMapWrapper.forEach(outlet.outletMap._outlets, (v, k) => this.unloadOutlet(v, components));
            if (this.performMutation) {
                outlet.unload();
            }
            else {
                this.deactivations.push(components.concat([outlet.loadedComponent]));
            }
        }
    }
}
//# sourceMappingURL=router.js.map