import {Component, AfterViewInit} from "@angular/core";
import {BlockSyncService} from "../services/skycoin.sync.service";
import {Observable} from "rxjs";
import {SyncProgress} from "../model/sync.progress";
declare var _: any;

@Component({
  selector: 'skycoin-block-sync',
  template: `
             <div *ngIf="syncDone == false">
               <h2>Syncing wallet, please wait .... </h2>
               <span *ngIf="currentWalletNumber >0">Current : {{currentWalletNumber}} out of {{highestWalletNumber}}</span>
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

