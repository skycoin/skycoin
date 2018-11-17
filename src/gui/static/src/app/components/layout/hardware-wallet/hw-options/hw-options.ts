import { Component, OnDestroy } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog';
import { HwSeedDialogComponent } from '../hw-seed-dialog/hw-seed-dialog';
import { ISubscription } from 'rxjs/Subscription';
import { WalletService } from '../../../../services/wallet.service';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog';
import { HwGenerateSeedDialogComponent } from '../hw-generate-seed-dialog/hw-generate-seed-dialog';
import { HwBackupDialogComponent } from '../hw-backup-dialog/hw-backup-dialog';

enum States {
  Disconnected,
  Processing,
  NewConnected,
  ConfiguredConnected,
  Error,
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
  private hwConnectionSubscription: ISubscription;

  private recheck = false;

  constructor(
    public dialogRef: MatDialogRef<HwWalletOptionsComponent>,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private walletService: WalletService,
  ) {
    this.checkWallet();
    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(() => this.checkWallet());
  }

  ngOnDestroy() {
    this.removeOperationSubscription();
    this.removeDialogSubscription();
    this.hwConnectionSubscription.unsubscribe();
  }

  generateMnemonic() {
    this.openDialog(HwGenerateSeedDialogComponent);
  }

  setMnemonic() {
    this.openDialog(HwSeedDialogComponent);
  }

  backup() {
    this.openDialog(HwBackupDialogComponent);
  }

  wipe() {
    this.openDialog(HwWipeDialogComponent);
  }

  private openDialog(dialogType, recheckAfterClosed = true) {
    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';

    if (recheckAfterClosed) {
      config.data = (() => this.recheck = true);
    }

    this.dialogSubscription = this.dialog.open(dialogType, config)
      .afterClosed().subscribe(() => {
        if (recheckAfterClosed && this.recheck) {
          this.recheck = false;
          this.checkWallet();
        }
      });
  }

  private removeDialogSubscription() {
    if (this.dialogSubscription) {
      this.dialogSubscription.unsubscribe();
    }
  }

  private removeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }

  private checkWallet() {
    if (!this.hwWalletService.getDeviceSync()) {
      this.currentState = States.Disconnected;
    } else {
      this.currentState = States.Processing;
      this.removeOperationSubscription();

      this. operationSubscription = this.hwWalletService.getAddresses(1, 0).subscribe(
        response => {
          this.walletAddress = response.rawResponse[0];
          this.walletService.wallets.first().subscribe(wallets => {
            const alreadySaved = wallets.some(wallet => wallet.addresses[0].address === this.walletAddress && wallet.isHardware);
            if (!alreadySaved) {
              this.walletService.createHardwareWallet(this.walletAddress);
              this.openDialog(HwAddedDialogComponent, false);
            }
            this.currentState = States.ConfiguredConnected;
          });
        },
        error => {
          if ((error as string).includes('Mnemonic not set')) {
            this.currentState = States.NewConnected;
          } else {
            this.currentState = States.Error;
          }
        },
      );
    }
  }
}
