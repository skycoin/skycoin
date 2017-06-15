System.register(['@angular/core', '@angular/http', 'rxjs/Observable', 'rxjs/add/operator/map', 'rxjs/add/operator/catch'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, http_1, http_2, Observable_1;
    var OutputService;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (Observable_1_1) {
                Observable_1 = Observable_1_1;
            },
            function (_1) {},
            function (_2) {}],
        execute: function() {
            OutputService = (function () {
                function OutputService(_http) {
                    this._http = _http;
                }
                OutputService.prototype.getOutPuts = function (addresses) {
                    var params = new http_1.URLSearchParams();
                    params.set('addrs', addresses.join());
                    return this._http
                        .get('/outputs', { headers: this.getHeaders(), search: params })
                        .map(function (res) { return res.json(); })
                        .catch(function (error) { return Observable_1.Observable.throw(error.json().error || 'Server error'); });
                };
                OutputService.prototype.getHeaders = function () {
                    var headers = new http_2.Headers();
                    headers.append('Accept', 'application/json');
                    return headers;
                };
                OutputService = __decorate([
                    core_1.Injectable(), 
                    __metadata('design:paramtypes', [http_1.Http])
                ], OutputService);
                return OutputService;
            }());
            exports_1("OutputService", OutputService);
        }
    }
});

//# sourceMappingURL=output.service.js.map
