import { OutputService } from '../services/output.service';
import { Wallet } from '../model/wallet.pojo';
import { Component, Input, AfterViewInit} from "@angular/core";
import {Output} from "../model/outputs.pojo";
declare var _: any;

@Component({
    selector: 'skycoin-outputs',
    template: `
                <div class="main-content ng-scope">
                            <table class="table table-bordered" style="width:100%">
                                <thead>
                                <tr class="dark-row">
                                    <th>Address</th>
                                    <th>Coins</th>
                                    <th>Hours</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr *ngFor="let item of outPuts" style="background:white">
                                    <td>{{item.address}}</td>
                                    <td>{{item.coins}}</td>
                                    <td>{{item.hours}}</td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
              `
    ,
    providers:[OutputService]
})

export class SkyCoinOutputComponent implements AfterViewInit {

    @Input()
    wallets:Wallet[];

    outPuts:Output[]=[];

    constructor(private _outputService:OutputService){}


    ngAfterViewInit(): any {
      this.refreshOutputs();
    }

    refreshOutputs():any{
        var addresses = _.flatten(this.wallets.map(function(item){return item.entries;})).map(function(item){return item.address});
        this._outputService.getOutPuts(addresses).subscribe(outputs=>
            {
                this.outPuts = outputs.head_outputs;
            },
            err => {
                console.log(err);
            }
        );
    }
}
