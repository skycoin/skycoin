import { Component, AfterViewInit} from "@angular/core";
import {PendingTransactionService} from "../services/pending.transaction.service";
declare var moment: any;

@Component({
  selector: 'pending-transactions',
  template: `
<button class="btn btn-default right" type="button" (click)="sendExecution()" >Resend for execution</button>
<div class="table-responsive">
                  <table id="pending-table" class="table">
                            <tbody>
                            <tr class="dark-row">
                                <td>S. No</td>
                                <td>Time received</td>
                                <td>Transaction ID</td>
                                <td>Inputs</td>
                                <td>Outputs</td>
                                <td>Amount</td>
                            </tr>
                            <tr *ngFor="let transaction of transactions;let i=index">
                                <td>{{i+1}}</td>
                                <td>{{getDateTimeString(transaction.received)}}</td>
                                <td>{{transaction.transaction.txid}}</td>
                                <td>
                                <p *ngFor="let input of transaction.transaction.inputs">{{input}},<br></p>
</td>
                                <td>
                                <p *ngFor="let output of transaction.transaction.outputs">{{output.dst}},<br></p>
</td>
                                <td><p *ngFor="let output of transaction.transaction.outputs">{{output.coins}},<br></p></td>
                            </tr>
                            </tbody>
                        </table>
                        </div>
              `
  ,
  providers:[PendingTransactionService]
})

export class PendingTxnsComponent implements AfterViewInit {


  constructor(private _pendingTxnsService:PendingTransactionService){}

  transactions:any[];

  ngAfterViewInit(): any {
    this.refreshPendingTxns();
  }

  refreshPendingTxns():any{
    this._pendingTxnsService.getPendingTransactions().subscribe(pendingTxns=>
        {
          this.transactions = pendingTxns;
        },
        err => {
          console.log(err);
        }
    );
  }
  getDateTimeString(ts) {
    return moment(ts).format("YYYY-MM-DD HH:mm")
  }
  sendExecution():void{
    this._pendingTxnsService.resendPendingTxns().subscribe(()=>{
      this.refreshPendingTxns();
    });
  }
}

