import {Component, OnInit, ViewChild} from 'angular2/core';
import { ROUTER_DIRECTIVES, OnActivate } from 'angular2/router';
import {Http, HTTP_BINDINGS, Response} from 'angular2/http';
import {HTTP_PROVIDERS, Headers} from 'angular2/http';
import {Observable} from 'rxjs/Observable';
import {Observer} from 'rxjs/Observer';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import {Dialog} from "./modal.ts";

@Component({
  selector: 'load-wallet',
  directives: [ROUTER_DIRECTIVES, Dialog],
  providers: [],
  templateUrl: 'app/templates/wallet.html'
})

export class loadWalletComponent implements OnInit {

  @ViewChild(Dialog) _dialog: Dialog;

  //public wallets;
  wallets : Array<any>;
  //balance: Array<any>;
  public progress;
  displayMode: DisplayModeEnum;
  displayModeEnum = DisplayModeEnum;

  constructor(private http: Http) { }

  ngOnInit() {
      this.displayMode = DisplayModeEnum.first;
      this.loadWallet();
      this.loadProgress();

      setInterval(() => {
        //this.loadWallet();
        console.log("Refreshing balance");
      }, 15000);
  }

  loadWallet(){
      this.http.post('/wallets', '')
        .map((res:Response) => res.json())
        .subscribe(
          data => {
            this.wallets = data;

            //Load Balance for each wallet
            var headers = new Headers();
            headers.append('Content-Type', 'application/x-www-form-urlencoded');
            var inc = 0;
            for(var item in data){
              var address = data[inc].meta.filename;
              this.http.post('/wallet/balance', JSON.stringify({id: address}), {headers: headers})
              .map((res:Response) => res.json())
              .subscribe(
                response => {
                  console.log('load done: ' + inc);
                  this.wallets[inc].balance =  response.confirmed.coins / 1000000;
                  inc++;
                },
                err => console.error("Error on load balance: "+err),
                () => console.log('Balance load done')
              );
            }
            //Load Balance for each wallet end

           },
          err => console.error("Error on load wallet: "+err),
          () => console.log('Wallet load done')
        );
  }

  loadProgress(){
      this.http.post('/blockchain/progress', '')
        .map((res:Response) => res.json())
        .subscribe(
          response => { this.progress = (parseInt(response.current,10)+1) / parseInt(response.Highest,10) * 100 },
          err => console.error("Error on load progress: "+err),
          () => console.log('Progress load done')
        );
  }

  switchTab(mode: DisplayModeEnum) {
      this.displayMode = mode;
  }
  
  showQR(wallet){
      this._dialog.showQR(wallet);
  }
  showNewWalletDialog(wallet){
      this._dialog.showWallet();
  }

}

enum DisplayModeEnum {
  first = 0,
  second = 1,
  third = 2
}
