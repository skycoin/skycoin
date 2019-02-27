import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-hw-passphrase-help-dialog',
  templateUrl: './hw-passphrase-help-dialog.component.html',
  styleUrls: ['./hw-passphrase-help-dialog.component.scss'],
})
export class HwPassphraseHelpDialogComponent {

  constructor(
    public dialogRef: MatDialogRef<HwPassphraseHelpDialogComponent>,
  ) { }
}
