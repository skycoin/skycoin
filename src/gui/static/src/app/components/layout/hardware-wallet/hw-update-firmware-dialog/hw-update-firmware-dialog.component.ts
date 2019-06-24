import { Component, ViewChild, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ButtonComponent } from '../../button/button.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { ISubscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';

enum States {
  Connecting,
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedTimeout,
}

@Component({
  selector: 'app-hw-update-firmware-dialog',
  templateUrl: './hw-update-firmware-dialog.component.html',
  styleUrls: ['./hw-update-firmware-dialog.component.scss'],
})
export class HwUpdateFirmwareDialogComponent extends HwDialogBaseComponent<HwUpdateFirmwareDialogComponent> implements OnDestroy {

  closeIfHwDisconnected = false;

  @ViewChild('button') button: ButtonComponent;

  currentState: States = States.Connecting;
  states = States;
  confirmed = false;

  deviceInBootloaderMode = false;
  deviceHasFirmware = true;

  private checkDeviceSubscription: ISubscription;

  get title(): string {
    if (this.currentState === States.Connecting) {
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
    this.currentState = States.Processing;

    this.closeCheckDeviceSubscription();

    this.operationSubscription = this.hwWalletService.updateFirmware().subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
      },
      err => {
        if (err.result !== null && err.result !== undefined && err.result === OperationResults.Success) {
          this.currentState = States.ReturnedSuccess;
        } else if (err.result && err.result === OperationResults.Timeout) {
          this.currentState = States.ReturnedTimeout;
        } else {
          if (err.result) {
            const errorMsg = getHardwareWalletErrorMsg(this.hwWalletService, this.translateService, err);
            setTimeout(() => {
              this.button.setError(errorMsg);
              this.msgBarService.showError(errorMsg);
            });
          } else {
            setTimeout(() => {
              this.button.setError(err);
              this.msgBarService.showError(err);
            });
          }

          this.checkDevice(false);

          this.currentState = States.Initial;
        }
      },
    );
  }

  private checkDevice(delay = true) {
    this.closeCheckDeviceSubscription();

    this.checkDeviceSubscription = Observable.of(0).delay(delay ? 1000 : 0).flatMap(() => this.hwWalletService.getFeatures(false)).subscribe(response => {
      this.deviceInBootloaderMode = response.rawResponse.bootloader_mode;
      if (this.deviceInBootloaderMode) {
        this.deviceHasFirmware = response.rawResponse.firmware_present;
      } else {
        this.deviceHasFirmware = true;
      }

      if (this.currentState === States.Connecting) {
        this.currentState = States.Initial;
      }

      this.checkDevice();
    }, () => {
      this.deviceInBootloaderMode = false;
      this.deviceHasFirmware = true;

      if (this.currentState === States.Connecting) {
        this.currentState = States.Initial;
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
