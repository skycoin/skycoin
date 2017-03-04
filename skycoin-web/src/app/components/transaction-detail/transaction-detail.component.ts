import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, Params} from "@angular/router";
import {Observable} from "rxjs";
import {TransactionDetailService} from "./transaction-detail.service";
import {Transaction} from "../block-chain-table/block";

@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.css']
})
export class TransactionDetailComponent implements OnInit {


  private transactionObservable:Observable<any>;

  private transaction:Transaction;


  constructor(   private service:TransactionDetailService,
                 private route: ActivatedRoute,
                 private router: Router) {
    this.transactionObservable=null;
    this.transaction = null;
  }

  ngOnInit() {
    this.transactionObservable= this.route.params
      .switchMap((params: Params) => {
        let txid = params['txid'];
        return this.service.getTransaction(txid);
      });

    this.transactionObservable.subscribe((trans)=>{
      this.transaction = trans;
      console.log(trans);
    })
  }

}
