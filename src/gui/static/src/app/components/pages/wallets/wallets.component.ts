import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet } from '../../../app.datatypes';
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
    private dialog: MatDialog,
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
    const config = new MatDialogConfig();
    config.width = '566px';
    this.dialog.open(CreateWalletComponent, config).afterClosed().subscribe();
  }

  loadWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    this.dialog.open(LoadWalletComponent, config).afterClosed().subscribe();
  }

  toggleWallet(wallet: Wallet) {
    wallet.opened = !wallet.opened;
  }
}
