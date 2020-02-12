import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog.component';
import { SubscriptionLike,  Observable } from 'rxjs';
import { HwAddedDialogComponent } from '../hw-added-dialog/hw-added-dialog.component';
import { HwGenerateSeedDialogComponent } from '../hw-generate-seed-dialog/hw-generate-seed-dialog.component';
import { HwBackupDialogComponent } from '../hw-backup-dialog/hw-backup-dialog.component';
import { HwChangePinDialogComponent } from '../hw-change-pin-dialog/hw-change-pin-dialog.component';
import { HwRestoreSeedDialogComponent } from '../hw-restore-seed-dialog/hw-restore-seed-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HwRemovePinDialogComponent } from '../hw-remove-pin-dialog/hw-remove-pin-dialog.component';
import { HwUpdateFirmwareDialogComponent } from '../hw-update-firmware-dialog/hw-update-firmware-dialog.component';
import { HwUpdateAlertDialogComponent } from '../hw-update-alert-dialog/hw-update-alert-dialog.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { map, first } from 'rxjs/operators';
import { AppConfig } from '../../../../app.config';
import { OperationError, HWOperationResults } from '../../../../utils/operation-error';
import { processServiceError } from '../../../../utils/errors';
import { HardwareWalletService, HwFeaturesResponse, HwSecurityWarnings } from '../../../../services/wallet-operations/hardware-wallet.service';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';

export interface ChildHwDialogParams {
  wallet: WalletBase;
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

  private dialogSubscription: SubscriptionLike;

  private completeRecheckRequested = false;
  private recheckSecurityOnlyRequested = false;
  private showErrorRequested = false;
  private wallet: WalletBase;

  public static openDialog(dialog: MatDialog, onboarding: boolean): MatDialogRef<HwOptionsDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = onboarding;
    config.autoFocus = false;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(HwOptionsDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) private onboarding: boolean,
    public dialogRef: MatDialogRef<HwOptionsDialogComponent>,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private hardwareWalletService: HardwareWalletService,
    private walletsAndAddressesService: WalletsAndAddressesService,
  ) {
    super(hwWalletService, dialogRef);

    this.checkWallet(true);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.removeDialogSubscription();
    this.removeOperationSubscription();
    this.msgBarService.hide();
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

    return this.hardwareWalletService.getFeaturesAndUpdateData(!dontUpdateWallet ? this.wallet : null).pipe(map((response: HwFeaturesResponse) => {
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

      if (response.walletNameUpdated) {
        this.msgBarService.showWarning('hardware-wallet.general.name-updated');
      }

      this.firmwareVersion = response.features.fw_major + '.' + response.features.fw_minor + '.' + response.features.fw_patch;

      return response;
    }));
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
    }, (err: OperationError) => {
      if (err.type && err.type === HWOperationResults.DaemonConnectionError) {
        this.showResult({
          text: 'hardware-wallet.errors.daemon-connection-with-configurable-link',
          link: AppConfig.hwWalletDaemonDownloadUrl,
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
        this.operationSubscription = this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
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
      (err: OperationError) => {
        err = processServiceError(err);

        if (err.type === HWOperationResults.Timeout) {
          this.operationSubscription = this.hwWalletService.getFeatures(false).subscribe(result => {
            if (result.rawResponse.bootloader_mode) {
              this.openUpdateDialog();
            } else {
              this.showError();
            }
          }, () => this.showError());
        } else if (err.type === HWOperationResults.WithoutSeed) {
          this.currentState = this.states.Finished;
          this.newWalletConnected = true;

          this.operationSubscription = this.updateSecurityWarningsAndData(true).subscribe(result => {
            if (suggestToUpdate && result.securityWarnings.find(warning => warning === HwSecurityWarnings.OutdatedFirmware)) {
              this.openUpdateWarning();
            }
          });
        } else if (err.type === HWOperationResults.FailedOrRefused) {
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = false;
        } else if (err.type === HWOperationResults.WrongPin) {
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = true;
        } else {
          this.processResult(err);
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
      text: this.customErrorMsg ? this.customErrorMsg : 'hardware-wallet.errors.generic-error',
      icon: this.msgIcons.Error,
    });
    this.customErrorMsg = '';
  }
}
