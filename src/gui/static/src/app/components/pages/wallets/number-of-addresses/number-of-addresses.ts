import { Component, OnInit, Inject, ViewChild, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, FormControl } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-number-of-addresses',
  templateUrl: './number-of-addresses.html',
  styleUrls: ['./number-of-addresses.css'],
})
export class NumberOfAddressesComponent implements OnInit, OnDestroy {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: any,
    public dialogRef: MatDialogRef<NumberOfAddressesComponent>,
    private msgBarService: MsgBarService,
  ) {}

  ngOnInit() {
    this.form = new FormGroup({});
    this.form.addControl('quantity', new FormControl(1, [this.validateQuantity]));
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

  private validateQuantity(control: FormControl) {
    if (control.value < 1 || control.value > 100 || Number(control.value) !== Math.round(Number(control.value))) {
      return { invalid: true };
    }

    return null;
  }
}
