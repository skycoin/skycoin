System.register(['../services/seed.service', '../model/seed.pojo', "@angular/core"], function(exports_1, context_1) {
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
    var seed_service_1, seed_pojo_1, core_1;
    var SeedComponent;
    return {
        setters:[
            function (seed_service_1_1) {
                seed_service_1 = seed_service_1_1;
            },
            function (seed_pojo_1_1) {
                seed_pojo_1 = seed_pojo_1_1;
            },
            function (core_1_1) {
                core_1 = core_1_1;
            }],
        execute: function() {
            SeedComponent = (function () {
                function SeedComponent(_seedService) {
                    this._seedService = _seedService;
                    this.seedValue = '';
                    this.currentSeed = new seed_pojo_1.Seed('');
                }
                SeedComponent.prototype.ngOnInit = function () {
                    var _this = this;
                    this._seedService.getMnemonicSeed().subscribe(function (seedReceived) {
                        _this.currentSeed = seedReceived;
                        _this.seedValue = seedReceived.seed;
                    }, function (err) {
                        console.log(err);
                    });
                };
                SeedComponent.prototype.getCurrentSeed = function () {
                    return this.seedValue;
                };
                SeedComponent = __decorate([
                    core_1.Component({
                        selector: 'seed-mnemonic',
                        template: "\n                 <textarea rows=\"4\"  placeholder=\"Wallet Seed\" cols=\"46\" class=\"form-control\" [(ngModel)]=\"seedValue\"></textarea>\n              ",
                        providers: [seed_service_1.SeedService]
                    }), 
                    __metadata('design:paramtypes', [seed_service_1.SeedService])
                ], SeedComponent);
                return SeedComponent;
            }());
            exports_1("SeedComponent", SeedComponent);
        }
    }
});

//# sourceMappingURL=seed.component.js.map
