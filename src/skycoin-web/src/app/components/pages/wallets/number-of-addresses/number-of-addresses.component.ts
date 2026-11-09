import { ChangeDetectionStrategy, Component } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup, ValidationErrors } from '@angular/forms';
import { MatDialogConfig, MatDialogRef } from '@angular/material/dialog';

import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

export const maxAddressesToAdd = 100;

/** Asks how many addresses to add to a wallet. Closes with the number, or nothing if cancelled. */
@Component({
    selector: 'app-number-of-addresses',
    templateUrl: './number-of-addresses.component.html',
    styleUrls: ['./number-of-addresses.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class NumberOfAddressesComponent {
  form = new UntypedFormGroup({
    quantity: new UntypedFormControl(1, [validQuantity]),
  });

  static openDialog(dialog: CustomMatDialogService): MatDialogRef<NumberOfAddressesComponent, number | undefined> {
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = true;

    return dialog.open(NumberOfAddressesComponent, config);
  }

  constructor(public dialogRef: MatDialogRef<NumberOfAddressesComponent>) { }

  closePopup() {
    this.dialogRef.close();
  }

  continue() {
    if (this.form.valid) {
      this.dialogRef.close(Number(this.form.value.quantity));
    }
  }
}

function validQuantity(control: AbstractControl): ValidationErrors | null {
  const value = Number(control.value);

  return Number.isInteger(value) && value >= 1 && value <= maxAddressesToAdd ? null : { quantity: true };
}
