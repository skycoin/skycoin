import {Component, View} from 'angular2/core';
//import {ROUTER_DIRECTIVES, RouteParams, RouteConfig} from 'angular2/router';
//import {RouteParams, RouterLink, RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS, LocationStrategy, HashLocationStrategy} from 'angular2/router';
import {RouterLink} from 'angular2/router';
import {bootstrap} from 'angular2/platform/browser';
import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
//import {loadSend} from "./app.send.ts";
import {loadWalletComponent} from "./app.loadWallet.ts";

export class initview {

}

bootstrap(loadWalletComponent,[HTTP_BINDINGS, HTTP_PROVIDERS]);
