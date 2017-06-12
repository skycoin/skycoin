import { Injectable } from '@angular/core';
import {Http, Response, Headers} from "@angular/http";
import {Observable} from "rxjs";
import {AddressBalance, Wallet} from "../home/models/models";

@Injectable()
export class WalletService {

  constructor(private _http: Http) {

  }
  private getHeaders() {
    var headers = new Headers();
    headers.append('Content-Type', 'application/x-www-form-urlencoded');
    headers.append('access-Control-Allow-Origin', '*');
    return headers;
  }

  getWallets():Observable<Wallet[]>{
    return this._http.post('http://127.0.0.1:6420/wallets', this.getHeaders())
    .map((res:Response) => res.json());
  }

  getCurrentBalanceOfAddress(address:string): Observable<AddressBalance> {
    return this._http.get('http://127.0.0.1:6420/balance?addrs='+address)
    .map((res:Response) => {
      return res.json()})
    .catch((error:any) => {
      console.log(error);
      return Observable.throw(error || 'Server error');
    });
  }

}
