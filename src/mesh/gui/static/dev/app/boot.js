System.register(['angular2/platform/browser', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', "./app.skywire.js"], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var browser_1, http_1, http_2, app_skywire_ts_1;
    return {
        setters:[
            function (browser_1_1) {
                browser_1 = browser_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (app_skywire_ts_1_1) {
                app_skywire_ts_1 = app_skywire_ts_1_1;
            }],
        execute: function() {
            browser_1.bootstrap(app_skywire_ts_1.loadComponent, [http_1.HTTP_BINDINGS, http_2.HTTP_PROVIDERS]);
        }
    }
});

//# sourceMappingURL=boot.js.map
