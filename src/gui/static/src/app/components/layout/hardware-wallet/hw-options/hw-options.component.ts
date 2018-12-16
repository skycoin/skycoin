import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog.component';
import { ISubscription } from 'rxjs/Subscription';
import { WalletService, SecurityWarnings } from '../../../../services/wallet.service';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog.component';
import { HwGenerateSeedDialogComponent } from '../hw-generate-seed-dialog/hw-generate-seed-dialog.component';
import { HwBackupDialogComponent } from '../hw-backup-dialog/hw-backup-dialog.component';
import { MessageIcons } from '../hw-message/hw-message.component';
import { Wallet } from '../../../../app.datatypes';
import { HwChangePinDialogComponent } from '../hw-change-pin-dialog/hw-change-pin-dialog.component';
import { HwRestoreSeedDialogComponent } from '../hw-restore-seed-dialog/hw-restore-seed-dialog.component';
import { Observable } from 'rxjs/Observable';

enum States {
  Disconnected,
  Processing,
  NewConnected,
  ConfiguredConnected,
  Error,
  ReturnedRefused,
  WrongPin,
}

@Component({
  selector: 'app-hw-options',
  templateUrl: './hw-options.component.html',
  styleUrls: ['./hw-options.component.scss'],
})
export class HwWalletOptionsComponent implements OnDestroy {

  msgIcons = MessageIcons;
  currentState: States;
  states = States;
  walletName = '';
  customErrorMsg = '';
  needsBackup: boolean;
  needsPin: boolean;

  private operationSubscription: ISubscription;
  private dialogSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  private completeRecheckRequested = false;
  private recheckSecurityOnlyRequested = false;
  private showErrorRequested = false;
  private wallet: Wallet;

  constructor(
    @Inject(MAT_DIALOG_DATA) private onboarding: boolean,
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
    this.openDialog(HwWipeDialogComponent, true);
  }

  wipeWithoutWallet() {
    this.openDialog(HwWipeDialogComponent);
  }

  closeModal() {
    this.dialogRef.close();
  }

  private openDialog(dialogType, includeWallet = false) {
    this.customErrorMsg = '';

    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    config.data = ((error: string = null, recheckSecurityOnly: boolean = false) => {
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
    });

    if (includeWallet) {
      config.data = {
        wallet: this.wallet,
        notifyFinishFunction: config.data,
      };
    }

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

  private removeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }

  private updateSecurityWarnings(): Observable<SecurityWarnings[]> {
    return this.walletService.updateWalletHasSecurityWarnings(this.wallet).map(warnings => {
      this.needsBackup = warnings.includes(SecurityWarnings.NeedsBackup);
      this.needsPin = warnings.includes(SecurityWarnings.NeedsPin);

      return warnings;
    });
  }

  private checkWallet() {
    if (!this.hwWalletService.getDeviceSync()) {
      this.currentState = States.Disconnected;
    } else {
      this.currentState = States.Processing;
      this.removeOperationSubscription();

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
          if (err.rawResponse && typeof err.rawResponse === 'string' && (err.rawResponse as string).includes('Mnemonic not set')) {
            this.currentState = States.NewConnected;
          } else {
            if (err.result && err.result === OperationResults.FailedOrRefused) {
              this.currentState = States.ReturnedRefused;
            } else if (err.result && err.result === OperationResults.WrongPin) {
              this.currentState = States.WrongPin;
            } else {
              this.currentState = States.Error;
            }
          }
        },
      );
    }
  }
}
