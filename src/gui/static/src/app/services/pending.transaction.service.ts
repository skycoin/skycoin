/**
 * Created by napandey on 4/6/17.
 */
import {Injectable} from '@angular/core';
import {Http, Response} from '@angular/http';
import {Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {Seed} from '../model/seed.pojo';
import {PendingTxn} from "../model/PendingTransaction";

@Injectable()
export class PendingTransactionService {
  constructor(private _http: Http) {
  }
  getPendingTransactions(): Observable<PendingTxn[]> {
    return this._http
    .get('/pendingTxs', {headers: this.getHeaders()})
    .map((res:Response) => res.json())
    .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
  }

  resendPendingTxns(): any {
    return this._http
    .get('/resendUnconfirmedTxns', {headers: this.getHeaders()})
    .map((res:Response) => res.json())
    .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
  }

  private getHeaders() {
    let headers = new Headers();
    headers.append('Accept', 'application/json');
    return headers;
  }
}

