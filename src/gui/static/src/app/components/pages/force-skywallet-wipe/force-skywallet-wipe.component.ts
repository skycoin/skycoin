import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { ChildHwDialogParams } from '../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { HwWipeDialogComponent } from '../../layout/hardware-wallet/hw-wipe-dialog/hw-wipe-dialog.component';

/**
 * Allows wipe a Skywallet. It is mainly for restoring the device when the user is not able to
 * use it because the operations are always cancelled for inactivity.
 */
@Component({
  selector: 'app-force-skywallet-wipe',
  templateUrl: './force-skywallet-wipe.component.html',
  styleUrls: ['./force-skywallet-wipe.component.scss'],
})
export class ForceSkywalletWipeComponent {
  constructor(
    public formBuilder: FormBuilder,
    private router: Router,
    private dialog: MatDialog,
  ) { }

  proceed() {
    const config = new MatDialogConfig();
    config.width = '450px';
    config.autoFocus = false;

    // Data for the modal window.
    config.data = <ChildHwDialogParams> {
      wallet: null,
      requestOptionsComponentRefresh: ((error: string = null, recheckSecurityOnly: boolean = false) => {
        if (!error) {
          // Return to the wallet list after the operation is done.
          this.router.navigate([''], {replaceUrl: true});
        }
      }),
    };

    this.dialog.open(HwWipeDialogComponent, config);
  }

  back() {
    // Return to the wallet list.
    this.router.navigate(['']);
  }
}
