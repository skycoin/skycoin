import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog.component';
import { ISubscription } from 'rxjs/Subscription';
import { WalletService, HwSecurityWarnings, HwFeaturesResponse } from '../../../../services/wallet.service';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog.component';
import { HwGenerateSeedDialogComponent } from '../hw-generate-seed-dialog/hw-generate-seed-dialog.component';
import { HwBackupDialogComponent } from '../hw-backup-dialog/hw-backup-dialog.component';
import { Wallet } from '../../../../app.datatypes';
import { HwChangePinDialogComponent } from '../hw-change-pin-dialog/hw-change-pin-dialog.component';
import { HwRestoreSeedDialogComponent } from '../hw-restore-seed-dialog/hw-restore-seed-dialog.component';
import { Observable } from 'rxjs/Observable';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HwWalletDaemonService } from '../../../../services/hw-wallet-daemon.service';
import { HwRemovePinDialogComponent } from '../hw-remove-pin-dialog/hw-remove-pin-dialog.component';
import { HwUpdateFirmwareDialogComponent } from '../hw-update-firmware-dialog/hw-update-firmware-dialog.component';
import { HwUpdateAlertDialogComponent } from '../hw-update-alert-dialog/hw-update-alert-dialog.component';

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

  newWalletConnected = false;
  otherStateBecauseWrongPin = false;
  walletName = '';
  customErrorMsg = '';
  firmwareVersion = '';

  securityWarnings: string[] = [];
  firmwareVersionNotVerified: boolean;
  outdatedFirmware: boolean;
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

    this.checkWallet(true);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.removeDialogSubscription();
    this.removeOperationSubscription();
  }

  hwConnectionChanged(connected: boolean) {
    this.checkWallet(true);
  }

  update() {
    this.openUpdateDialog();
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

  removePin() {
    this.openDialog(HwRemovePinDialogComponent);
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
          this.updateSecurityWarningsAndData(false, true).subscribe();
        } else if (this.showErrorRequested) {
          this.showError();
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

  private updateSecurityWarningsAndData(dontUpdateWallet = false, waitForResetingCurrentWarnings = false): Observable<HwFeaturesResponse> {
    if (!waitForResetingCurrentWarnings) {
      this.securityWarnings = [];
    }

    return this.walletService.getHwFeaturesAndUpdateData(!dontUpdateWallet ? this.wallet : null).map(response => {
      if (waitForResetingCurrentWarnings) {
        this.securityWarnings = [];
      }

      if (response.securityWarnings.includes(HwSecurityWarnings.FirmwareVersionNotVerified)) {
        this.firmwareVersionNotVerified = true;
        this.securityWarnings.push('hardware-wallet.options.unchecked-version-warning');
      } else {
        this.firmwareVersionNotVerified = false;
      }

      if (response.securityWarnings.includes(HwSecurityWarnings.OutdatedFirmware)) {
        this.outdatedFirmware = true;
        this.securityWarnings.push('hardware-wallet.options.outdated-version-warning');
      } else {
        this.outdatedFirmware = false;
      }

      if (!dontUpdateWallet && response.securityWarnings.includes(HwSecurityWarnings.NeedsBackup)) {
        this.needsBackup = true;
        this.securityWarnings.push('hardware-wallet.options.backup-warning');
      } else {
        this.needsBackup = false;
      }

      if (!dontUpdateWallet && response.securityWarnings.includes(HwSecurityWarnings.NeedsPin)) {
        this.needsPin = true;
        this.securityWarnings.push('hardware-wallet.options.pin-warning');
      } else {
        this.needsPin = false;
      }

      if (!dontUpdateWallet) {
        this.walletName = this.wallet.label;
      }

      this.firmwareVersion = response.features.fw_major + '.' + response.features.fw_minor + '.' + response.features.fw_patch;

      return response;
    });
  }

  private checkWallet(suggestToUpdate = false) {
    this.wallet = null;
    this.showResult({
      text: 'hardware-wallet.options.connecting',
      icon: this.msgIcons.Spinner,
    }, false);

    this.removeOperationSubscription();

    this.operationSubscription = this.hwWalletService.getDeviceConnected().subscribe(connected => {
      if (!connected) {
        this.showResult({
          text: 'hardware-wallet.options.disconnected',
          icon: this.msgIcons.Usb,
        });
      } else {
        this.operationSubscription = this.hwWalletService.getFeatures(false).subscribe(result => {
          if (result.rawResponse.bootloader_mode) {
            this.openUpdateDialog();
          } else {
            this.continueCheckingWallet(suggestToUpdate);
          }
        }, () => this.continueCheckingWallet(suggestToUpdate));
      }
    }, err => {
      if (err['_body'] && err['_body'] === HwWalletDaemonService.errorConnectingWithTheDaemon) {
        this.showResult({
          text: 'hardware-wallet.errors.daemon-connection',
          icon: this.msgIcons.Error,
        });
      } else {
        this.showError();
      }
    });
  }

  private continueCheckingWallet(suggestToUpdate) {
    this.operationSubscription = this.hwWalletService.getAddresses(1, 0).subscribe(
      response => {
        this.operationSubscription = this.walletService.wallets.first().subscribe(wallets => {
          const alreadySaved = wallets.some(wallet => {
            const found = wallet.addresses[0].address === response.rawResponse[0] && wallet.isHardware;
            if (found) {
              this.wallet = wallet;
              this.walletName = wallet.label;
            }

            return found;
          });
          if (alreadySaved) {
            this.operationSubscription = this.updateSecurityWarningsAndData().subscribe(result => {
              if (suggestToUpdate && result.securityWarnings.find(warning => warning === HwSecurityWarnings.OutdatedFirmware)) {
                this.openUpdateWarning();
              }

              if (!this.onboarding) {
                this.currentState = this.states.Finished;
                this.newWalletConnected = false;
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
        if (err.result && err.result === OperationResults.Timeout) {
          this.operationSubscription = this.hwWalletService.getFeatures(false).subscribe(result => {
            if (result.rawResponse.bootloader_mode) {
              this.openUpdateDialog();
            } else {
              this.showError();
            }
          }, () => this.showError());
        } else if (err.result && err.result === OperationResults.WithoutSeed) {
          this.currentState = this.states.Finished;
          this.newWalletConnected = true;

          this.operationSubscription = this.updateSecurityWarningsAndData(true).subscribe(result => {
            if (suggestToUpdate && result.securityWarnings.find(warning => warning === HwSecurityWarnings.OutdatedFirmware)) {
              this.openUpdateWarning();
            }
          });
        } else if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = false;
        } else if (err.result && err.result === OperationResults.WrongPin) {
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = true;
        } else {
          this.processResult(err.result);
        }
      },
    );
  }

  private removeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }

  private openUpdateWarning() {
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    this.dialog.open(HwUpdateAlertDialogComponent, config).afterClosed().subscribe(update => {
      if (update) {
        this.openUpdateDialog();
      }
    });
  }

  private openUpdateDialog() {
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    this.dialog.open(HwUpdateFirmwareDialogComponent, config);

    this.closeModal();
  }

  private showError() {
    this.showResult({
      text: this.customErrorMsg ? this.customErrorMsg : 'hardware-wallet.general.generic-error',
      icon: this.msgIcons.Error,
    });
    this.customErrorMsg = '';
  }
}
