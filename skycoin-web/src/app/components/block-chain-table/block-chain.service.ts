import { Injectable } from '@angular/core';
import {Http, Response} from '@angular/http';
import {Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {Block, BlockResponse} from "./block";

@Injectable()
export class BlockChainService {

  constructor(private _http: Http) { }

  getBlocks(startNumber:number, endNumber:number): Observable<Block[]> {
    var stringConvert = 'start='+startNumber+'&end='+endNumber;
    return this._http.get('/api/blocks?'+stringConvert)
      .map((res:Response) => {
      console.log(res);
      return res.json()})
      .map((res:BlockResponse)=>res.blocks)
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }

  getBlockByHash(hashNumber:string): Observable<Block> {
    var stringConvert = 'hash='+hashNumber;
    return this._http.get('/api/block?'+stringConvert)
      .map((res:Response) => {
        console.log(res);
        return res.json()})
      .map((res:Block)=>res)
      .catch((error:any) => {
        console.log(error);
        return Observable.throw(error || 'Server error');
      });
  }




}
