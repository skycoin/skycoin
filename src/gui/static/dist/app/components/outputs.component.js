System.register(['../services/output.service', "@angular/core"], function(exports_1, context_1) {
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
    var output_service_1, core_1;
    var SkyCoinOutputComponent;
    return {
        setters:[
            function (output_service_1_1) {
                output_service_1 = output_service_1_1;
            },
            function (core_1_1) {
                core_1 = core_1_1;
            }],
        execute: function() {
            SkyCoinOutputComponent = (function () {
                function SkyCoinOutputComponent(_outputService) {
                    this._outputService = _outputService;
                    this.outPuts = [];
                }
                SkyCoinOutputComponent.prototype.ngAfterViewInit = function () {
                    this.refreshOutputs();
                };
                SkyCoinOutputComponent.prototype.refreshOutputs = function () {
                    var _this = this;
                    var addresses = _.flatten(this.wallets.map(function (item) { return item.entries; })).map(function (item) { return item.address; });
                    this._outputService.getOutPuts(addresses).subscribe(function (outputs) {
                        _this.outPuts = outputs.head_outputs;
                    }, function (err) {
                        console.log(err);
                    });
                };
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Array)
                ], SkyCoinOutputComponent.prototype, "wallets", void 0);
                SkyCoinOutputComponent = __decorate([
                    core_1.Component({
                        selector: 'skycoin-outputs',
                        template: "\n                <div class=\"main-content ng-scope\">\n                            <table class=\"table table-bordered\" style=\"width:100%\">\n                                <thead>\n                                <tr class=\"dark-row\">\n                                    <th>Address</th>\n                                    <th>Coins</th>\n                                    <th>Hours</th>\n                                </tr>\n                                </thead>\n                                <tbody>\n                                <tr *ngFor=\"let item of outPuts\" style=\"background:white\">\n                                    <td>{{item.address}}</td>\n                                    <td>{{item.coins}}</td>\n                                    <td>{{item.hours}}</td>\n                                </tr>\n                                </tbody>\n                            </table>\n                        </div>\n              ",
                        providers: [output_service_1.OutputService]
                    }), 
                    __metadata('design:paramtypes', [output_service_1.OutputService])
                ], SkyCoinOutputComponent);
                return SkyCoinOutputComponent;
            }());
            exports_1("SkyCoinOutputComponent", SkyCoinOutputComponent);
        }
    }
});

//# sourceMappingURL=outputs.component.js.map
