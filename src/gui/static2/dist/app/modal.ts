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
export class Dialog{
    //Constructor method to load HTTP object
    constructor(private http: Http) { }

    //Declare default varialbes
    wallets : Array<any>;
    QrAddress: string;
    QrIsVisible: boolean;
    NewWalletIsVisible: boolean;
    EditWalletIsVisible: boolean;
    walletname: string;
    walletId: string;

    //Show QR code function for view QR popup
    showQR(address: any){
        this.QrAddress = address.entries[0].address;
        this.QrIsVisible = true;
    }
    //Hide QR code function for hide QR popup
    hideQr(){
        this.QrIsVisible = false;
    }

    //Show wallet function for view New wallet popup
    showWallet(){
        this.NewWalletIsVisible = true;
    }
    //Hide wallet function for hide New wallet popup
    hideWallet(){
        this.NewWalletIsVisible = false;
    }

    //Show edit wallet function
    showEditWallet(wallet: any){
        this.EditWalletIsVisible = true;
        this.walletId = wallet.meta.filename;
    }
    //Hide edit wallet function
    hideEditWallet(){
        this.EditWalletIsVisible = false;
    }

    //Add new wallet function for generate new wallet in Skycoin
    generateWallet(){
      //Set http headers
      var headers = new Headers();
      headers.append('Content-Type', 'application/x-www-form-urlencoded');

      //Post method executed
      this.http.post('/wallet/create', JSON.stringify({name: ''}), {headers: headers})
      .map((res:Response) => res.json())
      .subscribe(
        response => {
          //Load all wallets after creating new wallet
          this.http.post('/wallets', '')
            .map((res:Response) => res.json())
            .subscribe(
              //Response from API
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
