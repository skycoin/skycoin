import { Injectable } from '@angular/core';
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs";
import {UnspentOutput,AddressBalanceResponse} from "./UnspentOutput";
import {Block, Transaction} from "../block-chain-table/block";

@Injectable()
export class UxOutputsService {



  constructor(private _http: Http) { }

  getUxOutputsForAddress(address:number): Observable<UnspentOutput[]> {
    return this._http.get('/api/address?address='+address)
      .map((res:Response) => {
        return res.json()})
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }

  getCurrentBalanceOfAddress(address:number): Observable<AddressBalanceResponse> {
    return this._http.get('/api/currentBalance?address='+address)
    .map((res:Response) => {
      return res.json()})
    .catch((error:any) => {
      console.log(error);
      return Observable.throw(error || 'Server error');
    });
  }


  getAddressFromUxId(uxid:string): Observable<UnspentOutput[]> {
    return this._http.get('/api/uxout?uxid='+uxid)
      .map((res:Response) => {
        return res.json()})
      .map((res:any) =>{
        return res.owner_address
      })

      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }

  getBlockSource(blockNumber:number): Observable<Block> {
    return this._http.get('/api/blocks?start='+blockNumber+'&end='+blockNumber)
      .map((res:Response) => {
        return res.json().blocks[0]})
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }

}
