import { Component, Inject, OnInit } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';
import { UntypedFormBuilder } from '@angular/forms';

import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { copyTextToClipboard } from '../../../../../utils/general-utils';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { AppConfig } from '../../../../../app.config';

/**
 * Settings for CopyRawTxComponent.
 */
export interface CopyRawTxData {
  /**
   * Raw transaction text.
   */
  rawTx: string;
  /**
   * If the transaction is signed or not.
   */
  isUnsigned: boolean;
}

/**
 * Allows to see and copy a raw transaction, which is a hex string.
 */
@Component({
  selector: 'app-copy-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class CopyRawTxComponent extends OfflineDialogsBaseComponent implements OnInit {
  // Set the contents of some of the UI elements.
  cancelButtonText = 'common.close-button';
  okButtonText = 'offline-transactions.copy-tx.copy-button';
  inputLabel = 'offline-transactions.copy-tx.input-label';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, data: CopyRawTxData): MatDialogRef<CopyRawTxComponent, any> {
    const config = new MatDialogConfig();
    config.data = data;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(CopyRawTxComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: CopyRawTxData,
    public dialogRef: MatDialogRef<CopyRawTxComponent>,
    private msgBarService: MsgBarService,
    formBuilder: UntypedFormBuilder,
  ) {
    super(formBuilder);

    this.title = 'offline-transactions.copy-tx.' + (data.isUnsigned ? 'unsigned' : 'signed') + '-title';
    this.text = 'offline-transactions.copy-tx.text-' + (data.isUnsigned ? 'unsigned' : 'signed');
    this.contents = data.rawTx;

    this.currentState = OfflineDialogsStates.ShowingForm;
  }

  ngOnInit() {
    setTimeout(() => {
      this.okButton.focus();
    });
  }

  cancelPressed() {
    this.dialogRef.close();
  }

  okPressed() {
    // Copy the tx to the clipboad.
    copyTextToClipboard(this.data.rawTx);
    this.msgBarService.showDone('common.copied', 4000);
  }
}
