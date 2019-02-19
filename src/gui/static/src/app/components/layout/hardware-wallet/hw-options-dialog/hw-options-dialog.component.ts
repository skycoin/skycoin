import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog.component';
import { ISubscription } from 'rxjs/Subscription';
import { WalletService, HwSecurityWarnings } from '../../../../services/wallet.service';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog.component';
import { HwGenerateSeedDialogComponent } from '../hw-generate-seed-dialog/hw-generate-seed-dialog.component';
import { HwBackupDialogComponent } from '../hw-backup-dialog/hw-backup-dialog.component';
import { Wallet } from '../../../../app.datatypes';
import { HwChangePinDialogComponent } from '../hw-change-pin-dialog/hw-change-pin-dialog.component';
import { HwRestoreSeedDialogComponent } from '../hw-restore-seed-dialog/hw-restore-seed-dialog.component';
import { Observable } from 'rxjs/Observable';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

enum States {
  Disconnected,
  Processing,
  NewConnected,
  ConfiguredConnected,
  Error,
  ReturnedRefused,
  WrongPin,
}

export interface ChildHwDialogParams {
  wallet: Wallet;
  walletHasPin: boolean;
  requestOptionsComponentRefresh: any;
}

@Component({
  selector: 'app-hw-options-dialog',
  templateUrl: './hw-options-dialog.component.html',
  styleUrls: ['./hw-options-dialog.component.scss'],
})
export class HwOptionsDialogComponent extends HwDialogBaseComponent<HwOptionsDialogComponent> implements OnDestroy {

  closeIfHwDisconnected = false;

  currentState: States;
  states = States;
  walletName = '';
  customErrorMsg = '';
  needsBackup: boolean;
  needsPin: boolean;

  private dialogSubscription: ISubscription;

  private completeRecheckRequested = false;
  private recheckSecurityOnlyRequested = false;
  private showErrorRequested = false;
  private wallet: Wallet;

  constructor(
    @Inject(MAT_DIALOG_DATA) private onboarding: boolean,
    public dialogRef: MatDialogRef<HwOptionsDialogComponent>,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private walletService: WalletService,
  ) {
    super(hwWalletService, dialogRef);

    this.checkWallet();
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.removeDialogSubscription();
  }

  hwConnectionChanged(connected: boolean) {
    this.checkWallet();
  }

  generateMnemonic() {
    this.openDialog(HwGenerateSeedDialogComponent);
  }

  restoreMnemonic() {
    this.openDialog(HwRestoreSeedDialogComponent);
  }

  changePin() {
    this.openDialog(HwChangePinDialogComponent);
  }

  backup() {
    this.openDialog(HwBackupDialogComponent);
  }

  wipe() {
    this.openDialog(HwWipeDialogComponent);
  }

  confirmSeed() {
    this.openDialog(HwRestoreSeedDialogComponent);
  }

  private openDialog(dialogType) {
    this.customErrorMsg = '';

    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    config.data = <ChildHwDialogParams> {
      wallet: this.wallet,
      walletHasPin: !this.needsPin,
      requestOptionsComponentRefresh: ((error: string = null, recheckSecurityOnly: boolean = false) => {
        if (!error) {
          if (!recheckSecurityOnly) {
            this.completeRecheckRequested = true;
          } else {
            this.recheckSecurityOnlyRequested = true;
          }
        } else {
          this.showErrorRequested = true;
          this.customErrorMsg = error;
        }
      }),
    };

    this.dialogSubscription = this.dialog.open(dialogType, config)
      .afterClosed().subscribe(() => {
        if (this.completeRecheckRequested) {
          this.checkWallet();
        } else if (this.recheckSecurityOnlyRequested) {
          this.updateSecurityWarnings().subscribe();
        } else if (this.showErrorRequested) {
          this.currentState = States.Error;
        }
        this.completeRecheckRequested = false;
        this.recheckSecurityOnlyRequested = false;
        this.showErrorRequested = false;
      });
  }

  private removeDialogSubscription() {
    if (this.dialogSubscription) {
      this.dialogSubscription.unsubscribe();
    }
  }

  private updateSecurityWarnings(): Observable<HwSecurityWarnings[]> {
    return this.walletService.updateWalletHasHwSecurityWarnings(this.wallet).map(warnings => {
      this.needsBackup = warnings.includes(HwSecurityWarnings.NeedsBackup);
      this.needsPin = warnings.includes(HwSecurityWarnings.NeedsPin);

      return warnings;
    });
  }

  private checkWallet() {
    this.wallet = null;

    if (!this.hwWalletService.getDeviceConnectedSync()) {
      this.currentState = States.Disconnected;
    } else {
      this.currentState = States.Processing;

      if (this.operationSubscription) {
        this.operationSubscription.unsubscribe();
      }

      this.operationSubscription = this.hwWalletService.getAddresses(1, 0).subscribe(
        response => {
          this.walletService.wallets.first().subscribe(wallets => {
            const alreadySaved = wallets.some(wallet => {
              const found = wallet.addresses[0].address === response.rawResponse[0] && wallet.isHardware;
              if (found) {
                this.wallet = wallet;
                this.walletName = wallet.label;
              }

              return found;
            });
            if (alreadySaved) {
              this.updateSecurityWarnings().subscribe(() => {
                if (!this.onboarding) {
                  this.currentState = States.ConfiguredConnected;
                } else {
                  this.hwWalletService.showOptionsWhenPossible = true;
                  this.dialogRef.close(true);
                }
              });
            } else {
              this.openDialog(HwAddedDialogComponent);
            }
          });
        },
        err => {
          if (err.result && err.result === OperationResults.WithoutSeed) {
            this.currentState = States.NewConnected;
          } else if (err.result && err.result === OperationResults.FailedOrRefused) {
            this.currentState = States.ReturnedRefused;
          } else if (err.result && err.result === OperationResults.WrongPin) {
            this.currentState = States.WrongPin;
          } else {
            this.currentState = States.Error;
          }
        },
      );
    }
  }
}
