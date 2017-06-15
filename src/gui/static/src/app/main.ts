/* Avoid: 'error TS2304: Cannot find name <type>' during compilation */
///<reference path="../../typings/index.d.ts"/>

import {bootstrap} from "@angular/platform-browser-dynamic";
import {provide} from "@angular/core";
import {LocationStrategy, HashLocationStrategy} from "@angular/common";
import {ROUTER_PROVIDERS} from "@angular/router-deprecated";
import {HTTP_BINDINGS} from '@angular/http';
import {LoadWalletComponent} from "./app.loadWallet";

bootstrap(LoadWalletComponent, [
    ROUTER_PROVIDERS,
    HTTP_BINDINGS,
    provide(LocationStrategy, {useClass: HashLocationStrategy})
]);
