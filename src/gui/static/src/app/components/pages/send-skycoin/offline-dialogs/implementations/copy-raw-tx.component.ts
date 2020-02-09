import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { OfflineDialogsBaseComponent, OfflineDialogsStates } from '../offline-dialogs-base.component';
import { copyTextToClipboard } from '../../../../../utils/general-utils';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { FormBuilder } from '@angular/forms';
import { AppConfig } from '../../../../../app.config';

export interface CopyRawTxData {
  rawTx: string;
  isUnsigned: boolean;
}

@Component({
  selector: 'app-copy-raw-tx',
  templateUrl: '../offline-dialogs-base.component.html',
  styleUrls: ['../offline-dialogs-base.component.scss'],
})
export class CopyRawTxComponent extends OfflineDialogsBaseComponent implements OnInit {
  cancelButtonText = 'common.close-button';
  okButtonText = 'offline-transactions.copy-tx.copy-button';
  inputLabel = 'offline-transactions.copy-tx.input-label';

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
    formBuilder: FormBuilder,
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
    copyTextToClipboard(this.data.rawTx);
    this.msgBarService.showDone('common.copied', 4000);
  }
}
