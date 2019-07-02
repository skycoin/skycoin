import { Component, Inject, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { WalletService } from '../../../../services/wallet.service';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Wallet } from '../../../../app.datatypes';
import { ChangeNameComponent, ChangeNameData } from '../../../pages/wallets/change-name/change-name.component';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.component.html',
  styleUrls: ['./hw-added-dialog.component.scss'],
})
export class HwAddedDialogComponent extends HwDialogBaseComponent<HwAddedDialogComponent> implements OnDestroy {
  @ViewChild('input') input: ElementRef;
  wallet: Wallet;
  form: FormGroup;

  private initialLabel: string;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    private walletService: WalletService,
    hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.walletService.createHardwareWallet().subscribe(wallet => {
      this.operationSubscription = this.walletService.getHwFeaturesAndUpdateData(wallet).subscribe(() => {
        this.wallet = wallet;
        this.initialLabel = wallet.label;

        this.form = this.formBuilder.group({
          label: [wallet.label, Validators.required],
        });

        this.closeIfHwDisconnected = false;
        this.currentState = this.states.Finished;
        this.data.requestOptionsComponentRefresh();

        setTimeout(() => this.input.nativeElement.focus());
      }, err => this.processError(err));
    }, err => this.processError(err));
  }

  private processError(err: any) {
    if (err.result && err.result === OperationResults.Disconnected) {
      this.closeModal();

      return;
    }

    let errorMsg = 'hardware-wallet.general.generic-error-internet';

    if (err['_body']) {
      errorMsg = err['_body'];
    }
    this.showResult({
      text: errorMsg,
      icon: this.msgIcons.Error,
    });
    this.data.requestOptionsComponentRefresh(errorMsg);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.msgBarService.hide();
  }

  saveNameAndCloseModal() {
    if (this.form.value.label === this.initialLabel) {
      this.closeModal();
    } else {
      this.msgBarService.hide();

      const config = new MatDialogConfig();
      config.width = '400px';
      config.data = new ChangeNameData();
      (config.data as ChangeNameData).wallet = this.wallet;
      (config.data as ChangeNameData).newName = this.form.value.label;
      this.dialog.open(ChangeNameComponent, config).afterClosed().subscribe(result => {
        if (result && !result.errorMsg) {
          this.closeModal();
        } else if (result.errorMsg) {
          this.msgBarService.showError(result.errorMsg);
        }
      });
    }
  }
}
