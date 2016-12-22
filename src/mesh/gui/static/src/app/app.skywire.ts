import {Component, OnInit, ViewChild} from '@angular/core';
import {ROUTER_DIRECTIVES, OnActivate} from '@angular/router';

import {Http, HTTP_BINDINGS, Response} from '@angular/http';
import {HTTP_PROVIDERS, Headers} from '@angular/http';
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

export class LoadComponent implements OnInit {
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
              if (data && data.result && data.result.success) {
                self.nodes = data.orders||[];
              } else {
                return;
              }
            }, err => console.log("Error on load nodes: " + err), ()=>{});
    }
}
