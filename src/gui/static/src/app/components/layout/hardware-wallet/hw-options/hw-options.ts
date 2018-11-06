import { Component, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog';
import { HwSeedDialogComponent } from '../hw-seed-dialog/hw-seed-dialog';
import { ISubscription } from 'rxjs/Subscription';
import { WalletService } from '../../../../services/wallet.service';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog';

export enum States {
  Disconnected,
  Processing,
  NewConnected,
  ConfiguredConnected,
}

@Component({
  selector: 'app-hw-options',
  templateUrl: './hw-options.html',
  styleUrls: ['./hw-options.scss'],
})
export class HwWalletOptionsComponent implements OnDestroy {

  currentState: States;
  states = States;
  walletAddress = '';

  private operationSubscription: ISubscription;
  private dialogSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwWalletOptionsComponent>,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private walletService: WalletService,
  ) {
    this.checkWallet();
  }

  ngOnDestroy() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
    this.removeDialogSubscription();
  }

  setMnemonic() {
    this.openDialog(HwSeedDialogComponent);
  }

  wipe() {
    this.openDialog(HwWipeDialogComponent);
  }

  private openDialog(dialogType, recheckAfterClosed = true) {
    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';
    this.dialogSubscription = this.dialog.open(dialogType, config)
      .afterClosed().subscribe(() => {
        if (recheckAfterClosed) {
          this.checkWallet();
        }
      });
  }

  private removeDialogSubscription() {
    if (this.dialogSubscription) {
      this.dialogSubscription.unsubscribe();
    }
  }

  private checkWallet() {
    if (!this.hwWalletService.getDevice()) {
      this.currentState = States.Disconnected;
    } else {
      this.currentState = States.Processing;

      this. operationSubscription = this.hwWalletService.getAddresses(1, 0).subscribe(
        arg => {
          this.walletAddress = arg[0];
          this.walletService.wallets.first().subscribe(wallets => {
            const alreadySaved = wallets.some(wallet => wallet.addresses[0].address === this.walletAddress && wallet.isHardware);
            if (!alreadySaved) {
              this.walletService.createHardwareWallet(this.walletAddress);
              this.openDialog(HwAddedDialogComponent, false);
            }
            this.currentState = States.ConfiguredConnected;
          });
        },
        () => {
          this.currentState = States.NewConnected;
        },
      );
    }
  }
}
