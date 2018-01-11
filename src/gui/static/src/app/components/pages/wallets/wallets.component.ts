import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet } from '../../../app.datatypes';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { LoadWalletComponent } from './load-wallet/load-wallet.component';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss']
})
export class WalletsComponent implements OnInit, OnDestroy {

  coins: number;
  hours: number;

  private walletSubscription: Subscription;

  constructor(
    public walletService: WalletService,
    private dialog: MdDialog,
    private router: Router,
  ) {}

  ngOnInit() {
    this.walletSubscription = this.walletService.all().subscribe(wallets => {
      this.coins = wallets.map(wallet => wallet.coins >= 0 ? wallet.coins : 0).reduce((a , b) => a + b, 0);
      this.hours = wallets.map(wallet => wallet.hours >= 0 ? wallet.hours : 0).reduce((a , b) => a + b, 0);
    })
  }

  ngOnDestroy() {
    this.walletSubscription.unsubscribe();
  }

  addWallet() {
    const config = new MdDialogConfig();
    config.width = '566px';
    this.dialog.open(CreateWalletComponent, config).afterClosed().subscribe(result => {
      //
    });
  }

  loadWallet() {
    const config = new MdDialogConfig();
    config.width = '566px';
    this.dialog.open(LoadWalletComponent, config).afterClosed().subscribe(result => {
      //
    });
  }

  toggleWallet(wallet: Wallet) {
    wallet.opened = !wallet.opened;
  }
}
