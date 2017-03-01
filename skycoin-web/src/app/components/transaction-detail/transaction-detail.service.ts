import { Injectable } from '@angular/core';
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs";
import {Transaction} from "../block-chain-table/block";

@Injectable()
export class TransactionDetailService {


  constructor(private _http: Http) { }

  getTransaction(txid:string): Observable<Transaction[]> {
    return this._http.get('/api/transaction?txid='+txid)
      .map((res:Response) => {
        return res.json()})
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }


}
