import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, Params} from "@angular/router";
import {Observable} from "rxjs";
import 'rxjs/add/observable/forkJoin';
import {TransactionDetailService} from "./transaction-detail.service";
import {Transaction} from "../block-chain-table/block";
import * as moment from 'moment';

@Component({
  selector: 'app-transaction-detail',
  templateUrl: './transaction-detail.component.html',
  styleUrls: ['./transaction-detail.component.css']
})
export class TransactionDetailComponent implements OnInit {


  private transactionObservable:Observable<any>;

  private transaction:any;


  constructor(   private service:TransactionDetailService,
                 private route: ActivatedRoute,
                 private router: Router) {
    this.transactionObservable=null;
    this.transaction = null;
  }

  ngOnInit() {
    this.transactionObservable= this.route.params
      .flatMap((params: Params) => {
        let txid = params['txid'];
        return this.service.getTransaction(txid);
      })
    .flatMap((trans:any)=>{
      var tasks$ = [];
      this.transaction = trans.txn;
      this.transaction.status =trans.status.confirmed;
      this.transaction.block_num =trans.status.block_seq;
      trans=trans.txn;
      for(var i=0;i<trans.inputs.length;i++){
        tasks$.push(this.getAddressOfInput(trans.inputs[i]));
      }
      return Observable.forkJoin(...tasks$);
    });

    this.transactionObservable.subscribe((trans)=>{

      for(var i=0;i<trans.length;i++){
        this.transaction.inputs[i] = trans[i].owner_address;
      }
    })
  }

  getAddressOfInput(uxid:string):Observable<any>{
    return this.service.getInputAddress(uxid);
  }

  getTime(time:number){
    return moment.unix(time).format();
  }

}
