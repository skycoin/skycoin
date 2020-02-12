import { Component, Inject, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { ChangeNameComponent, ChangeNameData } from '../../../pages/wallets/change-name/change-name.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { OperationError, HWOperationResults } from '../../../../utils/operation-error';
import { processServiceError } from '../../../../utils/errors';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';

@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.component.html',
  styleUrls: ['./hw-added-dialog.component.scss'],
})
export class HwAddedDialogComponent extends HwDialogBaseComponent<HwAddedDialogComponent> implements OnDestroy {
  @ViewChild('input', { static: false }) input: ElementRef;
  wallet: WalletBase;
  form: FormGroup;
  maxHwWalletLabelLength = HwWalletService.maxLabelLength;

  private initialLabel: string;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private hardwareWalletService: HardwareWalletService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.walletsAndAddressesService.createHardwareWallet().subscribe(wallet => {
      this.operationSubscription = this.hardwareWalletService.getFeaturesAndUpdateData(wallet).subscribe(() => {
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

  private processError(err: OperationError) {
    err = processServiceError(err);
    if (err.type === HWOperationResults.Disconnected) {
      this.closeModal();

      return;
    }

    this.showResult({
      text: err.translatableErrorMsg,
      icon: this.msgIcons.Error,
    });
    this.data.requestOptionsComponentRefresh(err.translatableErrorMsg);
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

      const data = new ChangeNameData();
      data.wallet = this.wallet;
      data.newName = this.form.value.label;
      ChangeNameComponent.openDialog(this.dialog, data, true).afterClosed().subscribe(result => {
        if (result && !result.errorMsg) {
          this.closeModal();
        } else if (result && result.errorMsg) {
          this.msgBarService.showError(result.errorMsg);
        }
      });
    }
  }
}
