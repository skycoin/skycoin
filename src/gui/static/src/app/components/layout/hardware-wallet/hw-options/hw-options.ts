import { Component, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog';
import { HwSeedDialogComponent } from '../hw-seed-dialog/hw-seed-dialog';
import { ISubscription } from 'rxjs/Subscription';

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
  ) {
    this.checkWallet();
  }

  ngOnDestroy() {
    this.operationSubscription.unsubscribe();
    this.removeDialogSubscription();
  }

  setMnemonic() {
    this.openDialog(HwSeedDialogComponent);
  }

  wipe() {
    this.openDialog(HwWipeDialogComponent);
  }

  private openDialog(dialogType) {
    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';
    this.dialog.open(dialogType, config)
      .afterClosed().subscribe(() => this.checkWallet());
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
          this.currentState = States.ConfiguredConnected;
        },
        () => {
          this.currentState = States.NewConnected;
        },
      );
    }
  }
}
