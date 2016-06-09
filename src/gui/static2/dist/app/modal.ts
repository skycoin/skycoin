import { Component, OnInit } from 'angular2/core';
import { ROUTER_DIRECTIVES, OnActivate } from 'angular2/router';
import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS, Headers} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

import {QRCodeComponent} from './ng2-qrcode.ts'

@Component({
    selector: 'app-dialog',
    templateUrl: 'app/templates/modal.html',
    directives: [QRCodeComponent],
})
export class Dialog
{
    wallets : Array<any>;

    constructor(private http: Http) { }

    //For QR code dialog
    private QrAddress: string;
    public QrIsVisible: boolean;

    //For New Wallet dialog
    public NewWalletIsVisible: boolean;

    showQR(address: any){
        this.QrAddress = address.entries[0].address;
        this.QrIsVisible = true;
    }
    hideQr(){
        this.QrIsVisible = false;
    }

    showWallet(){
        this.NewWalletIsVisible = true;
    }
    hideWallet(){
        this.NewWalletIsVisible = false;
    }

    generateWallet(){
      alert("Oke");
      var headers = new Headers();
      headers.append('Content-Type', 'application/x-www-form-urlencoded');
      this.http.post('/wallet/create', JSON.stringify({name: ''}), {headers: headers})
      .map((res:Response) => res.json())
      .subscribe(
        response => {
          //Load all wallets after creating new wallet
          this.http.post('/wallets', '')
            .map((res:Response) => res.json())
            .subscribe(
              data => {
                this.NewWalletIsVisible = false;
                this.wallets = data;
              },
             err => console.error("Error on load wallet: "+err),
             () => console.log('Wallet load done')
           );
        },
        err => console.error("Error on create new wallet: "+err),
        () => console.log('New wallet create done')
      );
    }

}
