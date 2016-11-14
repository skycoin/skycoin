System.register(['angular2/core', 'angular2/router', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch'], function(exports_1, context_1) {
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
    var core_1, router_1, http_1, http_2;
    var loadComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {}],
        execute: function() {
            let loadComponent = class loadComponent {
                constructor(http) {
                    this.http = http;
                }
                //Init function for load default value
                ngOnInit() {
                    this.nodes = [];
                    this.transports = [];
                    this.loadNodeList();
                }
                loadNodeList() {
                    var self = this;
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    var url = '/nodemanager/getlistnodes';
                    this.http.get(url, { headers: headers })
                        .map((res) => res.json())
                        .subscribe(data => {
                        console.log("get node list", url, data);
                        if (data.result.success) {
                            self.nodes = data.orders;
                        }
                        else {
                            return;
                        }
                    }, err => console.log("Error on load nodes: " + err), () => { });
                }
            };
            loadComponent = __decorate([
                core_1.Component({
                    selector: 'load-skywire',
                    directives: [router_1.ROUTER_DIRECTIVES],
                    providers: [],
                    templateUrl: 'app/templates/template.html'
                }), 
                __metadata('design:paramtypes', [http_1.Http])
            ], loadComponent);
            exports_1("loadComponent", loadComponent);
        }
    }
});

//# sourceMappingURL=app.skywire.js.map
