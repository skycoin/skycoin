import {Component, AfterViewInit, Input, OnDestroy} from "@angular/core";
import {WalletService} from "../services/wallet.service";

declare var moment: any;
@Component({
  selector: 'backup-wallets',
  template: `
<h2>Wallet Backup</h2>
<p> Wallet Directory: <b>{{walletFolder}}</b> </p>
<p>
    <b>BACKUP YOUR SEED. ON PAPER. IN A SAFE PLACE.</b> As long as you have your seed, you can recover your coins.
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

                                <td><a class="btn btn-success" href="javascript:void(0);" (click)="download($event,wallet)">{{wallet.meta.filename}}</a></td>
                                 <td>
                                  <a class="btn btn-default" *ngIf="!wallet?.showSeed"  (click)="showOrHideSeed(wallet)">Show Seed</a>
                                  <p *ngIf="wallet?.showSeed">{{wallet.meta.seed}}<a class="btn btn-default btn-margin" (click)="showOrHideSeed(wallet)">Hide Seed</a></p>
                                 </td>
                            </tr>
                            </tbody>
                        </table>
                        </div>
              `
  ,
  styles: [`
    .btn-margin {
      margin: 0 1rem;
    }
  `],
  providers:[WalletService]
})

export class WalletBackupPageComponent implements AfterViewInit, OnDestroy{


  constructor(private _service:WalletService){}

  walletFolder:string;

  @Input()
  wallets:any[];


  ngAfterViewInit(): any {
    this.getWalletFolder();
    this.walletFolder = "";
  }

  ngOnDestroy() {
    this.wallets.forEach(el => {
      el.showSeed = false;
    })
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

  download(ev:Event,wallet:any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    let blob: Blob = new Blob([JSON.stringify({"seed":wallet.meta.seed})], { type: 'application/json'});
    let link=document.createElement('a');
    link.href=window.URL.createObjectURL(blob);
    link['download']= wallet.meta.filename + '.json';
    link.click();
  }
  showOrHideSeed(wallet){
    wallet.showSeed = !wallet.showSeed;
  }
}