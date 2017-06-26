import { Injectable } from '@angular/core';
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs";
import {CoinSupply} from "../block-chain-table/block";

@Injectable()
export class CoinSupplyService {

  constructor(private _http: Http) { }
  getCoinSupply(): Observable<CoinSupply> {
    return this._http.get('/api/coinSupply')
      .map((res:Response) => {
        return res.json()})
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }

}
