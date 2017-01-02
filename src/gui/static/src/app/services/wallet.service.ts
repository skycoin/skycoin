/**
 * Created by nakul.pandey@gmail.com on 01/01/17.
 */
import {Injectable} from '@angular/core';
import {Http, Response} from '@angular/http';
import {Headers} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {Wallet} from '../model/wallet.pojo';
declare var toastr: any;

@Injectable()
export class WalletService {

    wallets:Wallet[];

    constructor(private _http: Http) {
    }

    private getHeaders() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        return headers;
    }
    updateWallet(walletData): Observable<any> {
        var stringConvert = 'label='+walletData.newText+'&id='+walletData.walletId;
        return this._http.post('/wallet/update', stringConvert, {headers: this.getHeaders()})
            .map((res:Response) => res.json())
            .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }
}

