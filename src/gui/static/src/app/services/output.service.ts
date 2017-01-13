import {Injectable} from '@angular/core';
import {Http, Response, URLSearchParams} from '@angular/http';
import {Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {OutputsResponse} from '../model/outputs.pojo';

@Injectable()
export class OutputService {

    constructor(private _http: Http) {
    }

    getOutPuts(addresses:string[]): Observable<OutputsResponse> {
        let params = new URLSearchParams();
        params.set('addrs', addresses.join());
        return this._http
            .get('/outputs', {headers: this.getHeaders(), search:params})
            .map((res:Response) => res.json())
            .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }

    private getHeaders() {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return headers;
    }
}
