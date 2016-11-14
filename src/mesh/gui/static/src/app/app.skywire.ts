//import {Component, OnInit, ViewChild} from 'app/angular2/core';
//import {ROUTER_DIRECTIVES, OnActivate} from 'app/angular2/router';
import {Component, OnInit, ViewChild} from 'angular2/core';
import {ROUTER_DIRECTIVES, OnActivate} from 'angular2/router';

import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS, Headers} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

declare var _: any;
declare var $: any;
declare var moment: any;

@Component({
    selector: 'load-skywire',
    directives: [ROUTER_DIRECTIVES],
    providers: [],
    templateUrl: 'app/templates/template.html'
})

export class loadComponent implements OnInit {
    //Declare default varialbes
    nodes: Array<any>;
    transports: Array<any>;

    constructor(private http: Http) { }

    //Init function for load default value
    ngOnInit() {
      this.nodes = [];
      this.transports = [];
      this.loadNodeList();
    }

    loadNodeList() {
        var self = this;
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/nodemanager/getlistnodes';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get node list", url, data);
              if (data.result.success) {
                self.nodes = data.orders;
              } else {
                return;
              }
            }, err => console.log("Error on load nodes: " + err), ()=>{});
    }
}
