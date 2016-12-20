/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { Observable } from 'rxjs/Observable';
import { from } from 'rxjs/observable/from';
import { of } from 'rxjs/observable/of';
import { _catch } from 'rxjs/operator/catch';
import { concatAll } from 'rxjs/operator/concatAll';
import { first } from 'rxjs/operator/first';
import { map } from 'rxjs/operator/map';
import { mergeMap } from 'rxjs/operator/mergeMap';
import { EmptyError } from 'rxjs/util/EmptyError';
import { LoadedRouterConfig } from './router_config_loader';
import { NavigationCancelingError, PRIMARY_OUTLET, defaultUrlMatcher } from './shared';
import { UrlSegmentGroup, UrlTree } from './url_tree';
import { andObservables, forEach, merge, waitForMap, wrapIntoObservable } from './utils/collection';
var NoMatch = (function () {
    function NoMatch(segmentGroup) {
        if (segmentGroup === void 0) { segmentGroup = null; }
        this.segmentGroup = segmentGroup;
    }
    return NoMatch;
}());
var AbsoluteRedirect = (function () {
    function AbsoluteRedirect(urlTree) {
        this.urlTree = urlTree;
    }
    return AbsoluteRedirect;
}());
function noMatch(segmentGroup) {
    return new Observable(function (obs) { return obs.error(new NoMatch(segmentGroup)); });
}
function absoluteRedirect(newTree) {
    return new Observable(function (obs) { return obs.error(new AbsoluteRedirect(newTree)); });
}
function namedOutletsRedirect(redirectTo) {
    return new Observable(function (obs) { return obs.error(new Error("Only absolute redirects can have named outlets. redirectTo: '" + redirectTo + "'")); });
}
function canLoadFails(route) {
    return new Observable(function (obs) { return obs.error(new NavigationCancelingError("Cannot load children because the guard of the route \"path: '" + route.path + "'\" returned false")); });
}
export function applyRedirects(injector, configLoader, urlSerializer, urlTree, config) {
    return new ApplyRedirects(injector, configLoader, urlSerializer, urlTree, config).apply();
}
var ApplyRedirects = (function () {
    function ApplyRedirects(injector, configLoader, urlSerializer, urlTree, config) {
        this.injector = injector;
        this.configLoader = configLoader;
        this.urlSerializer = urlSerializer;
        this.urlTree = urlTree;
        this.config = config;
        this.allowRedirects = true;
    }
    ApplyRedirects.prototype.apply = function () {
        var _this = this;
        var expanded$ = this.expandSegmentGroup(this.injector, this.config, this.urlTree.root, PRIMARY_OUTLET);
        var urlTrees$ = map.call(expanded$, function (rootSegmentGroup) { return _this.createUrlTree(rootSegmentGroup, _this.urlTree.queryParams, _this.urlTree.fragment); });
        return _catch.call(urlTrees$, function (e) {
            if (e instanceof AbsoluteRedirect) {
                // after an absolute redirect we do not apply any more redirects!
                _this.allowRedirects = false;
                // we need to run matching, so we can fetch all lazy-loaded modules
                return _this.match(e.urlTree);
            }
            else if (e instanceof NoMatch) {
                throw _this.noMatchError(e);
            }
            else {
                throw e;
            }
        });
    };
    ApplyRedirects.prototype.match = function (tree) {
        var _this = this;
        var expanded$ = this.expandSegmentGroup(this.injector, this.config, tree.root, PRIMARY_OUTLET);
        var mapped$ = map.call(expanded$, function (rootSegmentGroup) {
            return _this.createUrlTree(rootSegmentGroup, tree.queryParams, tree.fragment);
        });
        return _catch.call(mapped$, function (e) {
            if (e instanceof NoMatch) {
                throw _this.noMatchError(e);
            }
            else {
                throw e;
            }
        });
    };
    ApplyRedirects.prototype.noMatchError = function (e) {
        return new Error("Cannot match any routes. URL Segment: '" + e.segmentGroup + "'");
    };
    ApplyRedirects.prototype.createUrlTree = function (rootCandidate, queryParams, fragment) {
        var root = rootCandidate.segments.length > 0 ?
            new UrlSegmentGroup([], (_a = {}, _a[PRIMARY_OUTLET] = rootCandidate, _a)) :
            rootCandidate;
        return new UrlTree(root, queryParams, fragment);
        var _a;
    };
    ApplyRedirects.prototype.expandSegmentGroup = function (injector, routes, segmentGroup, outlet) {
        if (segmentGroup.segments.length === 0 && segmentGroup.hasChildren()) {
            return map.call(this.expandChildren(injector, routes, segmentGroup), function (children) { return new UrlSegmentGroup([], children); });
        }
        else {
            return this.expandSegment(injector, segmentGroup, routes, segmentGroup.segments, outlet, true);
        }
    };
    ApplyRedirects.prototype.expandChildren = function (injector, routes, segmentGroup) {
        var _this = this;
        return waitForMap(segmentGroup.children, function (childOutlet, child) { return _this.expandSegmentGroup(injector, routes, child, childOutlet); });
    };
    ApplyRedirects.prototype.expandSegment = function (injector, segmentGroup, routes, segments, outlet, allowRedirects) {
        var _this = this;
        var routes$ = of.apply(void 0, routes);
        var processedRoutes$ = map.call(routes$, function (r) {
            var expanded$ = _this.expandSegmentAgainstRoute(injector, segmentGroup, routes, r, segments, outlet, allowRedirects);
            return _catch.call(expanded$, function (e) {
                if (e instanceof NoMatch)
                    return of(null);
                else
                    throw e;
            });
        });
        var concattedProcessedRoutes$ = concatAll.call(processedRoutes$);
        var first$ = first.call(concattedProcessedRoutes$, function (s) { return !!s; });
        return _catch.call(first$, function (e, _) {
            if (e instanceof EmptyError) {
                if (_this.noLeftoversInUrl(segmentGroup, segments, outlet)) {
                    return of(new UrlSegmentGroup([], {}));
                }
                else {
                    throw new NoMatch(segmentGroup);
                }
            }
            else {
                throw e;
            }
        });
    };
    ApplyRedirects.prototype.noLeftoversInUrl = function (segmentGroup, segments, outlet) {
        return segments.length === 0 && !segmentGroup.children[outlet];
    };
    ApplyRedirects.prototype.expandSegmentAgainstRoute = function (injector, segmentGroup, routes, route, paths, outlet, allowRedirects) {
        if (getOutlet(route) !== outlet)
            return noMatch(segmentGroup);
        if (route.redirectTo !== undefined && !(allowRedirects && this.allowRedirects))
            return noMatch(segmentGroup);
        if (route.redirectTo === undefined) {
            return this.matchSegmentAgainstRoute(injector, segmentGroup, route, paths);
        }
        else {
            return this.expandSegmentAgainstRouteUsingRedirect(injector, segmentGroup, routes, route, paths, outlet);
        }
    };
    ApplyRedirects.prototype.expandSegmentAgainstRouteUsingRedirect = function (injector, segmentGroup, routes, route, segments, outlet) {
        if (route.path === '**') {
            return this.expandWildCardWithParamsAgainstRouteUsingRedirect(injector, routes, route, outlet);
        }
        else {
            return this.expandRegularSegmentAgainstRouteUsingRedirect(injector, segmentGroup, routes, route, segments, outlet);
        }
    };
    ApplyRedirects.prototype.expandWildCardWithParamsAgainstRouteUsingRedirect = function (injector, routes, route, outlet) {
        var _this = this;
        var newTree = this.applyRedirectCommands([], route.redirectTo, {});
        if (route.redirectTo.startsWith('/')) {
            return absoluteRedirect(newTree);
        }
        else {
            return mergeMap.call(this.lineralizeSegments(route, newTree), function (newSegments) {
                var group = new UrlSegmentGroup(newSegments, {});
                return _this.expandSegment(injector, group, routes, newSegments, outlet, false);
            });
        }
    };
    ApplyRedirects.prototype.expandRegularSegmentAgainstRouteUsingRedirect = function (injector, segmentGroup, routes, route, segments, outlet) {
        var _this = this;
        var _a = match(segmentGroup, route, segments), matched = _a.matched, consumedSegments = _a.consumedSegments, lastChild = _a.lastChild, positionalParamSegments = _a.positionalParamSegments;
        if (!matched)
            return noMatch(segmentGroup);
        var newTree = this.applyRedirectCommands(consumedSegments, route.redirectTo, positionalParamSegments);
        if (route.redirectTo.startsWith('/')) {
            return absoluteRedirect(newTree);
        }
        else {
            return mergeMap.call(this.lineralizeSegments(route, newTree), function (newSegments) {
                return _this.expandSegment(injector, segmentGroup, routes, newSegments.concat(segments.slice(lastChild)), outlet, false);
            });
        }
    };
    ApplyRedirects.prototype.matchSegmentAgainstRoute = function (injector, rawSegmentGroup, route, segments) {
        var _this = this;
        if (route.path === '**') {
            if (route.loadChildren) {
                return map.call(this.configLoader.load(injector, route.loadChildren), function (r) {
                    route._loadedConfig = r;
                    return of(new UrlSegmentGroup(segments, {}));
                });
            }
            else {
                return of(new UrlSegmentGroup(segments, {}));
            }
        }
        else {
            var _a = match(rawSegmentGroup, route, segments), matched = _a.matched, consumedSegments_1 = _a.consumedSegments, lastChild = _a.lastChild;
            if (!matched)
                return noMatch(rawSegmentGroup);
            var rawSlicedSegments_1 = segments.slice(lastChild);
            var childConfig$ = this.getChildConfig(injector, route);
            return mergeMap.call(childConfig$, function (routerConfig) {
                var childInjector = routerConfig.injector;
                var childConfig = routerConfig.routes;
                var _a = split(rawSegmentGroup, consumedSegments_1, rawSlicedSegments_1, childConfig), segmentGroup = _a.segmentGroup, slicedSegments = _a.slicedSegments;
                if (slicedSegments.length === 0 && segmentGroup.hasChildren()) {
                    var expanded$ = _this.expandChildren(childInjector, childConfig, segmentGroup);
                    return map.call(expanded$, function (children) { return new UrlSegmentGroup(consumedSegments_1, children); });
                }
                else if (childConfig.length === 0 && slicedSegments.length === 0) {
                    return of(new UrlSegmentGroup(consumedSegments_1, {}));
                }
                else {
                    var expanded$ = _this.expandSegment(childInjector, segmentGroup, childConfig, slicedSegments, PRIMARY_OUTLET, true);
                    return map.call(expanded$, function (cs) { return new UrlSegmentGroup(consumedSegments_1.concat(cs.segments), cs.children); });
                }
            });
        }
    };
    ApplyRedirects.prototype.getChildConfig = function (injector, route) {
        var _this = this;
        if (route.children) {
            return of(new LoadedRouterConfig(route.children, injector, null, null));
        }
        else if (route.loadChildren) {
            return mergeMap.call(runGuards(injector, route), function (shouldLoad) {
                if (shouldLoad) {
                    if (route._loadedConfig) {
                        return of(route._loadedConfig);
                    }
                    else {
                        return map.call(_this.configLoader.load(injector, route.loadChildren), function (r) {
                            route._loadedConfig = r;
                            return r;
                        });
                    }
                }
                else {
                    return canLoadFails(route);
                }
            });
        }
        else {
            return of(new LoadedRouterConfig([], injector, null, null));
        }
    };
    ApplyRedirects.prototype.lineralizeSegments = function (route, urlTree) {
        var res = [];
        var c = urlTree.root;
        while (true) {
            res = res.concat(c.segments);
            if (c.numberOfChildren === 0) {
                return of(res);
            }
            else if (c.numberOfChildren > 1 || !c.children[PRIMARY_OUTLET]) {
                return namedOutletsRedirect(route.redirectTo);
            }
            else {
                c = c.children[PRIMARY_OUTLET];
            }
        }
    };
    ApplyRedirects.prototype.applyRedirectCommands = function (segments, redirectTo, posParams) {
        var t = this.urlSerializer.parse(redirectTo);
        return this.applyRedirectCreatreUrlTree(redirectTo, this.urlSerializer.parse(redirectTo), segments, posParams);
    };
    ApplyRedirects.prototype.applyRedirectCreatreUrlTree = function (redirectTo, urlTree, segments, posParams) {
        var newRoot = this.createSegmentGroup(redirectTo, urlTree.root, segments, posParams);
        return new UrlTree(newRoot, this.createQueryParams(urlTree.queryParams, this.urlTree.queryParams), urlTree.fragment);
    };
    ApplyRedirects.prototype.createQueryParams = function (redirectToParams, actualParams) {
        var res = {};
        forEach(redirectToParams, function (v, k) {
            if (v.startsWith(':')) {
                res[k] = actualParams[v.substring(1)];
            }
            else {
                res[k] = v;
            }
        });
        return res;
    };
    ApplyRedirects.prototype.createSegmentGroup = function (redirectTo, group, segments, posParams) {
        var _this = this;
        var updatedSegments = this.createSegments(redirectTo, group.segments, segments, posParams);
        var children = {};
        forEach(group.children, function (child, name) {
            children[name] = _this.createSegmentGroup(redirectTo, child, segments, posParams);
        });
        return new UrlSegmentGroup(updatedSegments, children);
    };
    ApplyRedirects.prototype.createSegments = function (redirectTo, redirectToSegments, actualSegments, posParams) {
        var _this = this;
        return redirectToSegments.map(function (s) { return s.path.startsWith(':') ? _this.findPosParam(redirectTo, s, posParams) :
            _this.findOrReturn(s, actualSegments); });
    };
    ApplyRedirects.prototype.findPosParam = function (redirectTo, redirectToUrlSegment, posParams) {
        var pos = posParams[redirectToUrlSegment.path.substring(1)];
        if (!pos)
            throw new Error("Cannot redirect to '" + redirectTo + "'. Cannot find '" + redirectToUrlSegment.path + "'.");
        return pos;
    };
    ApplyRedirects.prototype.findOrReturn = function (redirectToUrlSegment, actualSegments) {
        var idx = 0;
        for (var _i = 0, actualSegments_1 = actualSegments; _i < actualSegments_1.length; _i++) {
            var s = actualSegments_1[_i];
            if (s.path === redirectToUrlSegment.path) {
                actualSegments.splice(idx);
                return s;
            }
            idx++;
        }
        return redirectToUrlSegment;
    };
    return ApplyRedirects;
}());
function runGuards(injector, route) {
    var canLoad = route.canLoad;
    if (!canLoad || canLoad.length === 0)
        return of(true);
    var obs = map.call(from(canLoad), function (c) {
        var guard = injector.get(c);
        if (guard.canLoad) {
            return wrapIntoObservable(guard.canLoad(route));
        }
        else {
            return wrapIntoObservable(guard(route));
        }
    });
    return andObservables(obs);
}
function match(segmentGroup, route, segments) {
    var noMatch = { matched: false, consumedSegments: [], lastChild: 0, positionalParamSegments: {} };
    if (route.path === '') {
        if ((route.pathMatch === 'full') && (segmentGroup.hasChildren() || segments.length > 0)) {
            return { matched: false, consumedSegments: [], lastChild: 0, positionalParamSegments: {} };
        }
        else {
            return { matched: true, consumedSegments: [], lastChild: 0, positionalParamSegments: {} };
        }
    }
    var matcher = route.matcher || defaultUrlMatcher;
    var res = matcher(segments, segmentGroup, route);
    if (!res)
        return noMatch;
    return {
        matched: true,
        consumedSegments: res.consumed,
        lastChild: res.consumed.length,
        positionalParamSegments: res.posParams
    };
}
function split(segmentGroup, consumedSegments, slicedSegments, config) {
    if (slicedSegments.length > 0 &&
        containsEmptyPathRedirectsWithNamedOutlets(segmentGroup, slicedSegments, config)) {
        var s = new UrlSegmentGroup(consumedSegments, createChildrenForEmptySegments(config, new UrlSegmentGroup(slicedSegments, segmentGroup.children)));
        return { segmentGroup: mergeTrivialChildren(s), slicedSegments: [] };
    }
    else if (slicedSegments.length === 0 &&
        containsEmptyPathRedirects(segmentGroup, slicedSegments, config)) {
        var s = new UrlSegmentGroup(segmentGroup.segments, addEmptySegmentsToChildrenIfNeeded(segmentGroup, slicedSegments, config, segmentGroup.children));
        return { segmentGroup: mergeTrivialChildren(s), slicedSegments: slicedSegments };
    }
    else {
        return { segmentGroup: segmentGroup, slicedSegments: slicedSegments };
    }
}
function mergeTrivialChildren(s) {
    if (s.numberOfChildren === 1 && s.children[PRIMARY_OUTLET]) {
        var c = s.children[PRIMARY_OUTLET];
        return new UrlSegmentGroup(s.segments.concat(c.segments), c.children);
    }
    else {
        return s;
    }
}
function addEmptySegmentsToChildrenIfNeeded(segmentGroup, slicedSegments, routes, children) {
    var res = {};
    for (var _i = 0, routes_1 = routes; _i < routes_1.length; _i++) {
        var r = routes_1[_i];
        if (emptyPathRedirect(segmentGroup, slicedSegments, r) && !children[getOutlet(r)]) {
            res[getOutlet(r)] = new UrlSegmentGroup([], {});
        }
    }
    return merge(children, res);
}
function createChildrenForEmptySegments(routes, primarySegmentGroup) {
    var res = {};
    res[PRIMARY_OUTLET] = primarySegmentGroup;
    for (var _i = 0, routes_2 = routes; _i < routes_2.length; _i++) {
        var r = routes_2[_i];
        if (r.path === '' && getOutlet(r) !== PRIMARY_OUTLET) {
            res[getOutlet(r)] = new UrlSegmentGroup([], {});
        }
    }
    return res;
}
function containsEmptyPathRedirectsWithNamedOutlets(segmentGroup, slicedSegments, routes) {
    return routes
        .filter(function (r) { return emptyPathRedirect(segmentGroup, slicedSegments, r) &&
        getOutlet(r) !== PRIMARY_OUTLET; })
        .length > 0;
}
function containsEmptyPathRedirects(segmentGroup, slicedSegments, routes) {
    return routes.filter(function (r) { return emptyPathRedirect(segmentGroup, slicedSegments, r); }).length > 0;
}
function emptyPathRedirect(segmentGroup, slicedSegments, r) {
    if ((segmentGroup.hasChildren() || slicedSegments.length > 0) && r.pathMatch === 'full')
        return false;
    return r.path === '' && r.redirectTo !== undefined;
}
function getOutlet(route) {
    return route.outlet ? route.outlet : PRIMARY_OUTLET;
}
//# sourceMappingURL=apply_redirects.js.map