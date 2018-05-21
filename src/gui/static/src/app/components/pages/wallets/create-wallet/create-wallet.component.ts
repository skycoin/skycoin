import { Component, Inject, OnInit, ViewChild } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MAT_DIALOG_DATA } from '@angular/material';
import { ApiService } from '../../../../services/api.service';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss'],
})
export class CreateWalletComponent implements OnInit {
  @ViewChild('createButton') createButton: ButtonComponent;
  @ViewChild('cancelButton') cancelButton: ButtonComponent;
  form: FormGroup;
  seed: string;
  scan: Number;
  encrypt = true;
  disableDismiss = false;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data,
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
    private apiService: ApiService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    if (!this.form.valid || this.createButton.isLoading()) {
      return;
    }

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    const password = this.encrypt ? this.form.value.password : null;
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan, password)
      .subscribe(() => this.dialogRef.close(), e => {
        this.createButton.setError(e);
        this.cancelButton.disabled = false;
        this.disableDismiss = false;
      });
  }

  generateSeed(entropy: number) {
    this.apiService.generateSeed(entropy).subscribe(seed => this.form.get('seed').setValue(seed));
  }

  setEncrypt(event) {
    this.encrypt = event.checked;
    this.form.updateValueAndValidity();
  }

  private initForm() {
    this.form = new FormGroup({}, [this.validatePasswords.bind(this), this.validateSeeds.bind(this)]);
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('confirm_seed', new FormControl());
    this.form.addControl('password', new FormControl());
    this.form.addControl('confirm_password', new FormControl());

    if (this.data.create) {
      this.generateSeed(128);
    }

    this.scan = 100;
  }

  private validateSeeds() {
    if (this.data.create && this.form && this.form.get('seed') && this.form.get('confirm_seed')) {
      if (this.form.get('seed').value !== this.form.get('confirm_seed').value) {
        return { NotEqual: true };
      }
    }

    return null;
  }

  private validatePasswords() {
    if (this.encrypt && this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value) {
        if (this.form.get('password').value !== this.form.get('confirm_password').value) {
          return { NotEqual: true };
        }
      } else {
        return { Required: true };
      }
    }

    return null;
  }
}
