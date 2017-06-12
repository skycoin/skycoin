import {Injectable} from "@angular/core";
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs/Observable";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";
import {SyncProgress} from "../model/sync.progress";

@Injectable()
export class BlockSyncService {

  constructor(private _http: Http) {
  }

  getSyncProgress(): Observable<SyncProgress> {
    return this._http.post('/blockchain/progress', '')
    .map((res:Response) => res.json())
  }
}
