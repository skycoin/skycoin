import { Component, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MatDialogConfig, MatDialog, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { SubscriptionLike, Observable, of } from 'rxjs';
import { map, first, tap, mergeMap } from 'rxjs/operators';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwWipeDialogComponent } from '../hw-wipe-dialog/hw-wipe-dialog.component';
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
import { AppConfig } from '../../../../app.config';
import { OperationError, HWOperationResults } from '../../../../utils/operation-error';
import { processServiceError } from '../../../../utils/errors';
import { HardwareWalletService, HwFeaturesResponse, HwSecurityWarnings } from '../../../../services/wallet-operations/hardware-wallet.service';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';

/**
 * Params sent to the modal windows openned by HwOptionsDialogComponent.
 */
export interface ChildHwDialogParams {
  /**
   * Wallet object representing the connected device. It will be null if the connected device
   * have not been added to the wallet list.
   */
  wallet: WalletBase;
  /**
   * If the connected device has PIN code protection activated.
   */
  walletHasPin: boolean;
  /**
   * Function for asking the hw wallet options modal window to recheck the connected device,
   * due to changes made to it, and refresh the data.
   * @param error If a valid string is provided, instead of rechecking the connected device
   * the hw wallet options modal window will just hide the options and show the provided
   * error msg.
   * @param recheckSecurityOnly If true, instead of completelly recheck the connected device,
   * the hw wallet options modal window will just update the security warnings.
   */
  requestOptionsComponentRefresh(error?: string, recheckSecurityOnly?: boolean): void;
}

/**
 * Shows basic info and geenral configuration options for the connected hw wallet. If the
 * onboarding option is set to true, it will just try to find a known hw wallet connected
 * (if an unknown one is connected it will try to add it to the wallet list first and then
 * try to detect it again) and, after finding one, it will just close itself, return "true"
 * iin the "afterClosed" event and make a service wait to open it again after navigating
 * to the wallet list. This behavior allows the user to just add a wallet while on the wizard.
 */
@Component({
  selector: 'app-hw-options-dialog',
  templateUrl: './hw-options-dialog.component.html',
  styleUrls: ['./hw-options-dialog.component.scss'],
})
export class HwOptionsDialogComponent extends HwDialogBaseComponent<HwOptionsDialogComponent> implements OnDestroy {
  closeIfHwDisconnected = false;

  // If true, indicates that the connected device does not have a seed.
  newWalletConnected = false;
  // If true, the window state was set to "other" (which is used to show an error related to
  // the PIN code) because the user entered a wrong PIN. If false, it was because the user
  // did not enter the PIN code.
  otherStateBecauseWrongPin = false;
  // Version of the firmware of the connected device.
  firmwareVersion = '';

  // Wallet object corresponding to the current device, if it has been already added to the
  // wallet list.
  wallet: WalletBase;

  // Security warning texts to show on the UI.
  securityWarnings: string[] = [];
  // Vars for knowing which security warnings were found.
  firmwareVersionNotVerified: boolean;
  outdatedFirmware: boolean;
  needsBackup: boolean;
  needsPin: boolean;
  // Indicates if the security warning are being updated. The Security warnings are used as
  // a way to know the state of the device and show the appropiate options.
  refreshingWarnings = false;

  private dialogSubscription: SubscriptionLike;

  // If true, the last openned modal window requested, after closing it, this window to refesh
  // all its data.
  private completeRecheckRequested = false;
  // If true, the last openned modal window requested, after closing it, this window to refesh
  // the security warnings.
  private recheckSecurityOnlyRequested = false;
  // If true, the last openned modal window requested, after closing it, this window to show
  // a custom error msg.
  private showErrorRequested = false;
  // Custom error msg the last openned modal window requested to be shown.
  private customErrorMsg = '';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   * @param onboarding Set to true if the window is being openned from the wizard.
   */
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
    // If the device was connected or disconnected, recheck it.
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

  /**
   * Opens a modal window for making a specific operation related to the hw wallet.
   * @param dialogType Class of the modal window to open.
   */
  private openDialog(dialogType) {
    this.customErrorMsg = '';

    this.removeDialogSubscription();
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    // Data for the modal window.
    config.data = <ChildHwDialogParams> {
      // Include the current wallet, if there is one.
      wallet: this.wallet,
      walletHasPin: !this.needsPin,
      requestOptionsComponentRefresh: ((error: string = null, recheckSecurityOnly: boolean = false) => {
        if (!error) {
          // Set the data to be updated after closing the window, as requested.
          if (!recheckSecurityOnly) {
            this.completeRecheckRequested = true;
          } else {
            this.recheckSecurityOnlyRequested = true;
          }
        } else {
          // Set a custom error to be shown after closing the window.
          this.showErrorRequested = true;
          this.customErrorMsg = error;
        }
      }),
    };

    // Open the modal window.
    this.dialogSubscription = this.dialog.open(dialogType, config).afterClosed().subscribe(() => {
      // Refresh the data or show an error, if previously requested.
      if (this.completeRecheckRequested) {
        this.checkWallet();
      } else if (this.recheckSecurityOnlyRequested) {
        this.updateSecurityWarningsAndData().subscribe();
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

  /**
   * Updates the list of security warnings related to the connected device. Also updates and
   * saves the changes of the wallet object which represents the connected device, if there is
   * already one.
   * @returns The observable for performing the operation.
   */
  private updateSecurityWarningsAndData(): Observable<HwFeaturesResponse> {
    return of(1).pipe(
      tap(() => this.refreshingWarnings = true),
      mergeMap(() => this.hardwareWalletService.getFeaturesAndUpdateData(this.wallet)),
      map((response: HwFeaturesResponse) => {
        this.refreshingWarnings = false;
        this.securityWarnings = [];

        // Build a list with the texts of all warnings and set variables to known which
        // warnings were found.

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

        // If there is no wallet, the device does not have a seed, so the warning is not relevant.
        if (this.wallet && response.securityWarnings.includes(HwSecurityWarnings.NeedsBackup)) {
          this.needsBackup = true;
          this.securityWarnings.push('hardware-wallet.options.backup-warning');
        } else {
          this.needsBackup = false;
        }

        // If there is no wallet, the device does not have a seed, so the warning is not relevant.
        if (this.wallet && response.securityWarnings.includes(HwSecurityWarnings.NeedsPin)) {
          this.needsPin = true;
          this.securityWarnings.push('hardware-wallet.options.pin-warning');
        } else {
          this.needsPin = false;
        }

        // Inform if the name of the wallet on the wallet list was updated to match the one
        // shown on the device.
        if (response.walletNameUpdated) {
          this.msgBarService.showWarning('hardware-wallet.general.name-updated');
        }

        this.firmwareVersion = response.features.fw_major + '.' + response.features.fw_minor + '.' + response.features.fw_patch;

        return response;
      }),
    );
  }

  /**
   * Makes the initial steps for Checking if there is a connected hw wallet and its features.
   * @param suggestToUpdate If a modal window sugesting the user to update the firmware should
   * be shown if a more recent version of the firmware is detected.
   */
  private checkWallet(suggestToUpdate = false) {
    this.wallet = null;
    // Show the loading animation.
    this.showResult({
      text: 'hardware-wallet.options.connecting',
      icon: this.msgIcons.Spinner,
    }, false);

    this.removeOperationSubscription();

    this.operationSubscription = this.hwWalletService.getDeviceConnected().subscribe(connected => {
      if (!connected) {
        // Show the no device connected msg.
        this.showResult({
          text: 'hardware-wallet.options.disconnected',
          icon: this.msgIcons.Usb,
        });
      } else {
        // Try to check if the device is in bootloader mode, in which case the firmware update
        // modal window is openned. No previous operation is cancelled because that may cause
        // problems if the device is in bootloader mode. If the operation fails, it is asumed
        // that it was because the device has a pending operation and it is not in
        // bootloader mode.
        this.operationSubscription = this.hwWalletService.getFeatures(false).subscribe(result => {
          if (result.rawResponse.bootloader_mode) {
            this.openUpdateDialog();
          } else {
            this.continueCheckingWallet(suggestToUpdate);
          }
        }, () => this.continueCheckingWallet(suggestToUpdate));
      }
    }, (err: OperationError) => {
      this.processHwOperationError(err);
    });
  }

  /**
   * Checks if the connected device has a seed, updates its security warnings, verifies if
   * it was already added to the wallet list and any other operation nedded for showing the
   * appropiate info and options for the device.
   * @param suggestToUpdate If a modal window sugesting the user to update the firmware should
   * be shown if a more recent version of the firmware is detected.
   */
  private continueCheckingWallet(suggestToUpdate: boolean) {
    // Get the first address of the device, to use it for identification.
    this.operationSubscription = this.hwWalletService.getAddresses(1, 0).subscribe(
      response => {
        // If the first address was obteined, get all the saved wallets.
        this.operationSubscription = this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
          // Check if there is already a saved hw wallet with the obtained first address.
          const alreadySaved = wallets.some(wallet => {
            const found = wallet.addresses[0].address === response.rawResponse[0] && wallet.isHardware;
            if (found) {
              this.wallet = wallet;
            }

            return found;
          });

          if (alreadySaved) {
            // Update the security warnings.
            this.operationSubscription = this.updateSecurityWarningsAndData().subscribe(result => {
              if (suggestToUpdate && result.securityWarnings.find(warning => warning === HwSecurityWarnings.OutdatedFirmware)) {
                this.openUpdateWarning();
              }

              if (!this.onboarding) {
                // Show the options.
                this.currentState = this.states.Finished;
                this.newWalletConnected = false;
              } else {
                // Close the modal window and request to open it again after navigating to the
                // wallet list.
                this.hwWalletService.showOptionsWhenPossible = true;
                this.dialogRef.close(true);
              }
            }, err => this.processHwOperationError(err));
          } else {
            // Open the appropiate component for adding the device to the wallet list.
            this.openDialog(HwAddedDialogComponent);
          }
        });
      },
      (err: OperationError) => {
        err = processServiceError(err);

        if (err.type === HWOperationResults.WithoutSeed) {
          // If trying to get the first address failed because the device does not have a seed,
          // the options for a new device are shown.
          this.operationSubscription = this.updateSecurityWarningsAndData().subscribe(result => {
            this.currentState = this.states.Finished;
            this.newWalletConnected = true;

            if (suggestToUpdate && result.securityWarnings.find(warning => warning === HwSecurityWarnings.OutdatedFirmware)) {
              this.openUpdateWarning();
            }
          }, error => this.processHwOperationError(error));
        } else if (err.type === HWOperationResults.FailedOrRefused) {
          // Show an error due to a problem with the PIN and allow the user to wipe
          // the device if needed.
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = false;
        } else if (err.type === HWOperationResults.WrongPin) {
          // Show an error due to a problem with the PIN and allow the user to wipe
          // the device if needed.
          this.currentState = this.states.Other;
          this.otherStateBecauseWrongPin = true;
        } else {
          this.processHwOperationError(err);
        }
      },
    );
  }

  private removeOperationSubscription() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }

  /**
   * Opens a warning modal window asking the user to update the firmware.
   */
  private openUpdateWarning() {
    HwUpdateAlertDialogComponent.openDialog(this.dialog).afterClosed().subscribe(update => {
      if (update) {
        this.openUpdateDialog();
      }
    });
  }

  /**
   * Opens the firmware update modal window and closes this one.
   */
  private openUpdateDialog() {
    HwUpdateFirmwareDialogComponent.openDialog(this.dialog);
    this.closeModal();
  }

  /**
   * Shows the custom error msg saved on customErrorMsg, or a generic one, if no custom
   * msg was saved.
   */
  private showError() {
    this.showResult({
      text: this.customErrorMsg ? this.customErrorMsg : 'hardware-wallet.errors.generic-error',
      icon: this.msgIcons.Error,
    });
    this.customErrorMsg = '';
  }
}
