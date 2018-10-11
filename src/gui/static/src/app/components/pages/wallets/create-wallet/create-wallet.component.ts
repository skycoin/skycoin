import { Component, Inject, OnInit, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MAT_DIALOG_DATA } from '@angular/material';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss'],
})
export class CreateWalletComponent {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @ViewChild('createButton') createButton: ButtonComponent;
  @ViewChild('cancelButton') cancelButton: ButtonComponent;

  scan: Number;
  disableDismiss = false;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data,
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
  ) {}

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    if (!this.formControl.isValid || this.createButton.isLoading()) {
      return;
    }

    const data = this.formControl.getData();

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    this.walletService.create(data.label, data.seed, this.scan, data.password)
      .subscribe(() => this.dialogRef.close(), e => {
        this.createButton.setError(e);
        this.cancelButton.disabled = false;
        this.disableDismiss = false;
      });
  }
}
