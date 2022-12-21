import { Component, OnDestroy } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';
import { SubscriptionLike, of } from 'rxjs';
import { mergeMap, delay } from 'rxjs/operators';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { OperationError, HWOperationResults } from '../../../../utils/operation-error';
import { processServiceError } from '../../../../utils/errors';
import { AppConfig } from '../../../../app.config';

/**
 * Allows to download and install the firmware of the hw wallet.
 */
@Component({
  selector: 'app-hw-update-firmware-dialog',
  templateUrl: './hw-update-firmware-dialog.component.html',
  styleUrls: ['./hw-update-firmware-dialog.component.scss'],
})
export class HwUpdateFirmwareDialogComponent extends HwDialogBaseComponent<HwUpdateFirmwareDialogComponent> implements OnDestroy {
  closeIfHwDisconnected = false;

  currentState = this.states.Connecting;
  // If the user has confirmed the operation with the checkbox.
  confirmed = false;

  deviceInBootloaderMode = false;
  deviceHasFirmware = true;

  private checkDeviceSubscription: SubscriptionLike;

  // The texts shown on the modal window depend on the features of the connected device.

  get title(): string {
    if (this.currentState === this.states.Connecting) {
      return 'hardware-wallet.update-firmware.title-connecting';
    } else if (this.deviceHasFirmware) {
      return 'hardware-wallet.update-firmware.title-update';
    }

    return 'hardware-wallet.update-firmware.title-install';
  }

  get text(): string {
    if (!this.deviceHasFirmware) {
      return 'hardware-wallet.update-firmware.text-no-firmware';
    }

    if (this.deviceInBootloaderMode) {
      return 'hardware-wallet.update-firmware.text-bootloader';
    }

    return 'hardware-wallet.update-firmware.text-not-bootloader';
  }

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<HwUpdateFirmwareDialogComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.smallModalWidth;

    return dialog.open(HwUpdateFirmwareDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<HwUpdateFirmwareDialogComponent>,
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
  ) {
    super(hwWalletService, dialogRef);
    this.checkDevice(false);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.msgBarService.hide();
    this.closeCheckDeviceSubscription();
  }

  setConfirmed(event) {
    this.confirmed = event.checked;
  }

  // Downloads and installs the lastest firmware.
  startUpdating() {
    this.msgBarService.hide();
    this.showResult({
      text: 'hardware-wallet.update-firmware.text-downloading',
      icon: this.msgIcons.Spinner,
    });

    // Temporarily stop checking the device features.
    this.closeCheckDeviceSubscription();

    this.operationSubscription = this.hwWalletService.updateFirmware(() => this.currentState = this.states.Processing).subscribe(
      () => {
        // Show the success msg.
        this.showResult({
          text: 'hardware-wallet.update-firmware.finished',
          icon: this.msgIcons.Success,
        });
      },
      (err: OperationError) => {
        err = processServiceError(err);

        // Show the result.
        if (err.type === HWOperationResults.Success) {
          this.showResult({
            text: 'hardware-wallet.update-firmware.finished',
            icon: this.msgIcons.Success,
          });
        } else if (err.type === HWOperationResults.Timeout) {
          this.showResult({
            text: 'hardware-wallet.update-firmware.timeout',
            icon: this.msgIcons.Error,
          });
        } else {
          // If there was a simple error, return to the initial state.
          setTimeout(() => {
            this.msgBarService.showError(err);
          });

          this.checkDevice(false);

          this.currentState = this.states.Initial;
        }
      },
    );
  }

  /**
   * Checks if the device is connected and gets its features.
   * @param delayOperation If there must be a small delay before making the operation.
   */
  private checkDevice(delayOperation = true) {
    this.closeCheckDeviceSubscription();

    // The call to get the features asks no to cancel the current operation before getting
    // the data because the cancel operation does not work well in bootloader mode.
    this.checkDeviceSubscription = of(0).pipe(delay(delayOperation ? 1000 : 0), mergeMap(() => this.hwWalletService.getFeatures(false))).subscribe(response => {
      this.deviceInBootloaderMode = response.rawResponse.bootloader_mode;
      if (this.deviceInBootloaderMode) {
        this.deviceHasFirmware = response.rawResponse.firmware_present;
      } else {
        this.deviceHasFirmware = true;
      }

      if (this.currentState === this.states.Connecting) {
        this.currentState = this.states.Initial;
      }

      // Repeat the operation periodically.
      this.checkDevice();
    }, () => {
      // Asume the device is not connected.
      this.deviceInBootloaderMode = false;
      this.deviceHasFirmware = true;

      if (this.currentState === this.states.Connecting) {
        this.currentState = this.states.Initial;
      }

      // Repeat the operation periodically.
      this.checkDevice();
    });
  }

  private closeCheckDeviceSubscription() {
    if (this.checkDeviceSubscription) {
      this.checkDeviceSubscription.unsubscribe();
    }
  }
}
