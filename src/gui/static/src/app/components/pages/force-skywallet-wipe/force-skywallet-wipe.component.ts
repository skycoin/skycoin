import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { UntypedFormBuilder } from '@angular/forms';
import { MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { ChildHwDialogParams } from '../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { HwWipeDialogComponent } from '../../layout/hardware-wallet/hw-wipe-dialog/hw-wipe-dialog.component';
import { AppConfig } from '../../../app.config';

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
    public formBuilder: UntypedFormBuilder,
    private router: Router,
    private dialog: MatDialog,
  ) { }

  proceed() {
    const config = new MatDialogConfig();
    config.width = AppConfig.smallModalWidth;
    config.autoFocus = false;

    // Data for the modal window.
    config.data = {
      wallet: null,
      requestOptionsComponentRefresh: ((error: string = null, recheckSecurityOnly: boolean = false) => {
        if (!error) {
          // Return to the wallet list after the operation is done.
          this.router.navigate([''], {replaceUrl: true});
        }
      }),
    } as ChildHwDialogParams;

    this.dialog.open(HwWipeDialogComponent, config);
  }

  back() {
    // Return to the wallet list.
    this.router.navigate(['']);
  }
}
