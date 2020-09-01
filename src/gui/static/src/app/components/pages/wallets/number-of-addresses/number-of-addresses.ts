import { Component, OnInit, Inject, ViewChild, OnDestroy } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';

@Component({
  selector: 'app-number-of-addresses',
  templateUrl: './number-of-addresses.html',
  styleUrls: ['./number-of-addresses.scss'],
})
export class NumberOfAddressesComponent implements OnInit, OnDestroy {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;

  // Vars with the validation error messages.
  inputErrorMsg = '';

  public static openDialog(dialog: MatDialog, eventFunction: any): MatDialogRef<NumberOfAddressesComponent, any> {
    const config = new MatDialogConfig();
    config.data = eventFunction;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(NumberOfAddressesComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: any,
    public dialogRef: MatDialogRef<NumberOfAddressesComponent>,
    private msgBarService: MsgBarService,
  ) {}

  ngOnInit() {
    this.form = new FormGroup({});
    this.form.addControl('quantity', new FormControl(1));

    this.form.setValidators(this.validateForm.bind(this));
  }

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  closePopup() {
    this.dialogRef.close();
  }

  continue() {
    if (this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.button.setLoading();

    this.data(this.form.value.quantity, (close, endedWithError = false) => {
      this.button.resetState();
      if (!endedWithError) {
        if (close) {
          this.closePopup();
        }
      } else {
        this.msgBarService.showError('wallet.add-addresses.error');
      }
    });
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    // The number must be an integer from 1 to 100.
    const value = this.form.get('quantity').value as number;
    if (!value || value < 1 || value > 100 || value !== Math.round(value)) {
      valid = false;
      if (this.form.get('quantity').touched) {
        this.inputErrorMsg = 'wallet.add-addresses.quantity-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
