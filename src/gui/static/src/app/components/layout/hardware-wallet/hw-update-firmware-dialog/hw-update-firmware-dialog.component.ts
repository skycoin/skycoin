import { Component, ViewChild, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ButtonComponent } from '../../button/button.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';

enum States {
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

  currentState: States = States.Initial;
  states = States;
  confirmed = false;

  constructor(
    public dialogRef: MatDialogRef<HwUpdateFirmwareDialogComponent>,
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
  ) {
    super(hwWalletService, dialogRef);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.msgBarService.hide();
  }

  setConfirmed(event) {
    this.confirmed = event.checked;
  }

  startUpdating() {
    this.msgBarService.hide();
    this.currentState = States.Processing;

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

          this.currentState = States.Initial;
        }
      },
    );
  }
}
