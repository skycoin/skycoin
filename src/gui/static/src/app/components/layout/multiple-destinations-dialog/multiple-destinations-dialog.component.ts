import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MsgBarService } from '../../../services/msg-bar.service';
import { AppConfig } from '../../../app.config';

@Component({
  selector: 'app-multiple-destinations-dialog',
  templateUrl: './multiple-destinations-dialog.component.html',
  styleUrls: ['./multiple-destinations-dialog.component.scss'],
})
export class MultipleDestinationsDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;

  public static openDialog(dialog: MatDialog, content: string): MatDialogRef<MultipleDestinationsDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = content;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(MultipleDestinationsDialogComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<MultipleDestinationsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: string,
    private formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      data: [this.data],
    });
  }

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  processData() {
    try {
      if ((this.form.value.data as string).trim().length === 0) {
        this.msgBarService.showError('send.bulk-send.no-data-error');

        return;
      }

      let entries = (this.form.value.data as string).split(/\r?\n/);
      if (!entries || entries.length === 0) {
        this.msgBarService.showError('send.bulk-send.no-data-error');

        return;
      }

      entries = entries.filter(entry => entry.trim().length > 0);

      const firstElementParts = entries[0].split(',').length;
      if (firstElementParts !== 2 && firstElementParts !== 3) {
        this.msgBarService.showError('send.bulk-send.invalid-data-error');

        return;
      }

      const splitedEntries = [];
      let consistentNumberOfParts = true;
      entries.forEach((entry: string, i: number) => {
        splitedEntries[i] = entry.split(',');
        if (splitedEntries[i].length !== firstElementParts) {
          consistentNumberOfParts = false;
        }
      });

      if (!consistentNumberOfParts) {
        this.msgBarService.showError('send.bulk-send.inconsistent-data-error');

        return;
      }

      const response = [];
      splitedEntries.forEach((entry, i) => {
        response[i] = [];
        (entry as string[]).forEach((part, j) => {
          response[i][j] = part.trim();
        });
      });

      this.dialogRef.close(response);
    } catch (e) {
      this.msgBarService.showError('send.bulk-send.invalid-data-error');
    }
  }
}
