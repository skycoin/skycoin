import { Component, ViewChild, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ButtonComponent } from '../../button/button.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { SubscriptionLike, of } from 'rxjs';
import { mergeMap, delay } from 'rxjs/operators';

@Component({
  selector: 'app-hw-update-firmware-dialog',
  templateUrl: './hw-update-firmware-dialog.component.html',
  styleUrls: ['./hw-update-firmware-dialog.component.scss'],
})
export class HwUpdateFirmwareDialogComponent extends HwDialogBaseComponent<HwUpdateFirmwareDialogComponent> implements OnDestroy {

  closeIfHwDisconnected = false;

  @ViewChild('button', { static: false }) button: ButtonComponent;

  currentState = this.states.Connecting;
  confirmed = false;

  deviceInBootloaderMode = false;
  deviceHasFirmware = true;

  private checkDeviceSubscription: SubscriptionLike;

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

  constructor(
    public dialogRef: MatDialogRef<HwUpdateFirmwareDialogComponent>,
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
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

  startUpdating() {
    this.msgBarService.hide();
    this.showResult({
      text: 'hardware-wallet.update-firmware.text-downloading',
      icon: this.msgIcons.Spinner,
    });

    this.closeCheckDeviceSubscription();

    this.operationSubscription = this.hwWalletService.updateFirmware(() => this.currentState = this.states.Processing).subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.update-firmware.finished',
          icon: this.msgIcons.Success,
        });
      },
      err => {
        if (err.result !== null && err.result !== undefined && err.result === OperationResults.Success) {
          this.showResult({
            text: 'hardware-wallet.update-firmware.finished',
            icon: this.msgIcons.Success,
          });
        } else if (err.result && err.result === OperationResults.Timeout) {
          this.showResult({
            text: 'hardware-wallet.update-firmware.timeout',
            icon: this.msgIcons.Error,
          });
        } else {
          if (err.result) {
            const errorMsg = getHardwareWalletErrorMsg(this.translateService, err);
            setTimeout(() => {
              this.msgBarService.showError(errorMsg);
            });
          } else {
            setTimeout(() => {
              this.msgBarService.showError(err);
            });
          }

          this.checkDevice(false);

          this.currentState = this.states.Initial;
        }
      },
    );
  }

  private checkDevice(delayOperation = true) {
    this.closeCheckDeviceSubscription();

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

      this.checkDevice();
    }, () => {
      this.deviceInBootloaderMode = false;
      this.deviceHasFirmware = true;

      if (this.currentState === this.states.Connecting) {
        this.currentState = this.states.Initial;
      }

      this.checkDevice();
    });
  }

  private closeCheckDeviceSubscription() {
    if (this.checkDeviceSubscription) {
      this.checkDeviceSubscription.unsubscribe();
    }
  }
}
