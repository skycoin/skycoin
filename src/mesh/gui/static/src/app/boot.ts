import {Component} from 'angular2/core';
import {RouterLink} from 'angular2/router';
import {bootstrap} from 'angular2/platform/browser';
import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {loadComponent} from "./app.skywire.ts";

bootstrap(loadComponent,[HTTP_BINDINGS, HTTP_PROVIDERS]);
