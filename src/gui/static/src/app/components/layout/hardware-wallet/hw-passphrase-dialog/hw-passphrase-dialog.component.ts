import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HwPassphraseHelpDialogComponent } from '../hw-passphrase-help-dialog/hw-passphrase-help-dialog.component';

@Component({
  selector: 'app-hw-passphrase-dialog',
  templateUrl: './hw-passphrase-dialog.component.html',
  styleUrls: ['./hw-passphrase-dialog.component.scss'],
})
export class HwPassphraseDialogComponent extends HwDialogBaseComponent<HwPassphraseDialogComponent> implements OnInit, OnDestroy {
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<HwPassphraseDialogComponent>,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      passphrase: ['', Validators.required],
    });
  }

  ngOnDestroy() {
    super.ngOnDestroy();
  }

  sendPassphrase() {
    this.dialogRef.close(this.form.value.passphrase);
  }

  openHelp() {
    this.dialog.open(HwPassphraseHelpDialogComponent, <MatDialogConfig> {
      width: '450px',
    });
  }
}
