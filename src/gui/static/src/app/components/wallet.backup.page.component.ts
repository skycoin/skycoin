import {Component, AfterViewInit, Input} from "@angular/core";
import {WalletService} from "../services/wallet.service";
declare var moment: any;

@Component({
  selector: 'backup-wallets',
  template: `
<h2>Wallet Backup</h2>
<p> Your wallets are safely stored at : <b>{{walletFolder}}</b> </p>
<p>
    Store the wallet seed in a safe place. With this we can get back your wallet balance.
</p>
<div class="table-responsive">
                  <table id="wallet-table" class="table">
                  <thead>
                    <tr class="dark-row">
                                <td>S. No</td>
                                <td>Wallet Label</td>
                                <td>File Name</td>
                                <td>Download</td>
                                <td>Seed</td>
          
                            </tr>
</thead>
                            <tbody>
                      
                            <tr *ngFor="let wallet of wallets;let i=index">
                                <td>{{i+1}}</td>
                                <td>{{wallet.meta.label}}</td>
                                <td>{{wallet.meta.filename}}</td>

                                <td><a id="{{wallet.meta.seed}}" class="btn btn-success"  href="" download="{{getJsonObject(wallet)}}">{{wallet.meta.filename}}</a></td>
                                 <td id="seed-{{wallet.meta.seed}}"><a class="btn btn-default"  (click)="showSeed(wallet.meta.seed)">Show Seed</a></td>
                            </tr>
                            </tbody>
                        </table>
                        </div>
              `
  ,
  providers:[WalletService]
})

export class WalletBackupPageComponent implements AfterViewInit {


  constructor(private _service:WalletService){}

  walletFolder:string;

  @Input()
  wallets:any[];


  ngAfterViewInit(): any {
    this.getWalletFolder();
    this.walletFolder = "";
  }

  getWalletFolder():any{
    this._service.getWalletFolder().subscribe(walletFolder=>
        {
          this.walletFolder = walletFolder.address;
        },
        err => {
          console.log(err);
        }
    );
  }
  getJsonObject(wallet) {
    var walletDoc = document.getElementById(wallet.meta.seed);
    walletDoc.setAttribute("href","data:text/json;charset=utf-8," +encodeURIComponent(JSON.stringify({"seed":wallet.meta.seed})));
    return  wallet.meta.filename+'.json';
  }

  showSeed(seed){
    var seedEl = document.getElementById("seed-"+seed);
    seedEl.innerHTML = seed;
  }
}