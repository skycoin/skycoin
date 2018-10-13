import { Component, OnInit, Inject, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, FormControl } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';

@Component({
  selector: 'app-number-of-addresses',
  templateUrl: './number-of-addresses.html',
  styleUrls: ['./number-of-addresses.css'],
})
export class NumberOfAddressesComponent implements OnInit {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<NumberOfAddressesComponent>,
    private formBuilder: FormBuilder,
  ) {}

  ngOnInit() {
    this.form = new FormGroup({});
    this.form.addControl('quantity', new FormControl(1, [this.validateQuantity]));
  }

  closePopup() {
    this.dialogRef.close();
  }

  continue() {
    this.dialogRef.close(Math.round(Number(this.form.value.quantity)));
  }

  private validateQuantity(control: FormControl) {
    if (control.value < 1 || control.value > 100 || Number(control.value) !== Math.round(Number(control.value))) {
      return { invalid: true };
    }

    return null;
  }
}
