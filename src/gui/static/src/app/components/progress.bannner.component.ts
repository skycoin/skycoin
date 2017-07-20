import {Component, AfterViewInit} from "@angular/core";
import {BlockSyncService} from "../services/skycoin.sync.service";
import {Observable} from "rxjs";
import {SyncProgress} from "../model/sync.progress";
declare var _: any;

@Component({
  selector: 'skycoin-block-sync',
  template: `

        
             <div class="sync-div-container">
             
             <ul class="fa-ul">
  <li><i class="fa-li fa fa-spinner fa-spin" *ngIf="syncDone == false"></i>
  <span *ngIf="currentWalletNumber>0">{{currentWalletNumber}} of {{highestWalletNumber}} blocks synced</span>
  <span *ngIf="currentWalletNumber==0">Syncing wallet</span>
  </li>
</ul>
                
               
              </div>
            
              `
  ,
  providers:[BlockSyncService]
})

export class SkycoinSyncWalletBlock implements AfterViewInit {

  _syncProgress:Observable<SyncProgress>;

  syncDone:boolean;

  currentWalletNumber:number;
  highestWalletNumber:number;

  constructor(private _syncService:BlockSyncService){
    this.currentWalletNumber = this.highestWalletNumber =0;
    this.syncDone = false;
  }

  private handlerSync:any;

  ngAfterViewInit(): any {
    this.handlerSync = setInterval(() => {
          if(this.highestWalletNumber-this.currentWalletNumber<=1 && this.highestWalletNumber!=0){
            clearInterval(this.handlerSync);
            this.syncDone=true;
          }
          this.syncBlocks();
        }, 2000);
  }
  syncBlocks():any{
    this._syncProgress = this._syncService.getSyncProgress();
    this._syncProgress.subscribe((syncProgress:SyncProgress)=>{
      this.currentWalletNumber = syncProgress.current;
      this.highestWalletNumber = syncProgress.highest;
    });


  }
}

