import {Injectable} from '@angular/core';
import {Http, Response} from '@angular/http';
import {Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {Seed} from '../model/seed.pojo';

@Injectable()
export class SeedService {

    constructor(private _http: Http) {
    }

    getMnemonicSeed(): Observable<Seed> {
        return this._http
            .get('/wallet/newSeed', {headers: this.getHeaders()})
            .map((res:Response) => res.json())
            .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }

    private getHeaders() {
        let headers = new Headers();
        headers.append('Accept', 'application/json');
        return headers;
    }
}
