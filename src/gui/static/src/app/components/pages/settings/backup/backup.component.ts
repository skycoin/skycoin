import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { Wallet } from '../../../../app.datatypes';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SeedModalComponent } from './seed-modal/seed-modal.component';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';

@Component({
  selector: 'app-backup',
  templateUrl: './backup.component.html',
  styleUrls: ['./backup.component.scss'],
})
export class BackupComponent implements OnInit, OnDestroy {
  folder: string;
  wallets: Wallet[] = [];

  private walletSubscription;

  constructor(
    public walletService: WalletService,
    private dialog: MatDialog,
  ) {}

  ngOnInit() {
    this.walletService.folder().subscribe(folder => this.folder = folder);

    this.walletSubscription = this.walletService.all().subscribe(wallets => {
      this.wallets = wallets;
    });
  }

  ngOnDestroy() {
    this.walletSubscription.unsubscribe();
  }

  get onlyEncrypted() {
    return this.wallets.filter(wallet => wallet.encrypted);
  }

  showSeed(wallet: Wallet) {
    PasswordDialogComponent.openDialog(this.dialog, { wallet: wallet }).componentInstance.passwordSubmit
      .subscribe(passwordDialog => {
        this.walletService.getWalletSeed(wallet, passwordDialog.password).subscribe(seed => {
          passwordDialog.close();
          SeedModalComponent.openDialog(this.dialog, seed);
        }, err => passwordDialog.error(err));
      });
  }
}
