/* Avoid: 'error TS2304: Cannot find name <type>' during compilation */
///<reference path="../../typings/index.d.ts"/>
System.register(["@angular/platform-browser-dynamic", "@angular/core", "@angular/common", "@angular/router-deprecated", '@angular/http', "./app.loadWallet"], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var platform_browser_dynamic_1, core_1, common_1, router_deprecated_1, http_1, app_loadWallet_1;
    return {
        setters:[
            function (platform_browser_dynamic_1_1) {
                platform_browser_dynamic_1 = platform_browser_dynamic_1_1;
            },
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (common_1_1) {
                common_1 = common_1_1;
            },
            function (router_deprecated_1_1) {
                router_deprecated_1 = router_deprecated_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
            },
            function (app_loadWallet_1_1) {
                app_loadWallet_1 = app_loadWallet_1_1;
            }],
        execute: function() {
            platform_browser_dynamic_1.bootstrap(app_loadWallet_1.LoadWalletComponent, [
                router_deprecated_1.ROUTER_PROVIDERS,
                http_1.HTTP_BINDINGS,
                core_1.provide(common_1.LocationStrategy, { useClass: common_1.HashLocationStrategy })
            ]);
        }
    }
});

//# sourceMappingURL=main.js.map
