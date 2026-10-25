import { ChangeDetectionStrategy, Component, Inject } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs';

/** One destination as typed, before any validation of its values. */
export interface BulkDestination {
  address: string;
  coins: string;
  hours?: string;
}

/**
 * Reads one destination per line, as address,coins or address,coins,hours.
 * Either every line has hours or none does. Returns the destinations, or the
 * translation key of the error.
 */
export function parseBulkDestinations(text: string): BulkDestination[] | string {
  const lines = text.split(/\r?\n/).map(l => l.trim()).filter(l => l.length > 0);
  if (lines.length === 0) {
    return 'send.bulk-send.no-data-error';
  }

  const destinations: BulkDestination[] = [];
  let parts = 0;
  for (const line of lines) {
    const fields = line.split(',').map(f => f.trim());
    if ((fields.length !== 2 && fields.length !== 3) || fields.some(f => f.length === 0)) {
      return 'send.bulk-send.invalid-data-error';
    }
    if (parts !== 0 && fields.length !== parts) {
      return 'send.bulk-send.inconsistent-data-error';
    }
    parts = fields.length;
    destinations.push(fields.length === 3 ?
      { address: fields[0], coins: fields[1], hours: fields[2] } :
      { address: fields[0], coins: fields[1] });
  }

  return destinations;
}

/**
 * Takes several destinations as text. Closes with the destinations, or with
 * undefined when cancelled. Use BulkSendComponent.ask.
 */
@Component({
    selector: 'app-bulk-send',
    templateUrl: './bulk-send.component.html',
    styleUrls: ['./bulk-send.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class BulkSendComponent {
  form: UntypedFormGroup;
  error = '';

  static ask(dialog: MatDialog, current: string): Observable<BulkDestination[] | undefined> {
    return dialog.open(BulkSendComponent, { width: '600px', data: current }).afterClosed();
  }

  constructor(
    public dialogRef: MatDialogRef<BulkSendComponent>,
    @Inject(MAT_DIALOG_DATA) data: string,
    formBuilder: UntypedFormBuilder,
  ) {
    this.form = formBuilder.group({ data: [data] });
  }

  process() {
    const result = parseBulkDestinations(this.form.value.data || '');
    if (typeof result === 'string') {
      this.error = result;

      return;
    }
    this.dialogRef.close(result);
  }

  cancel() {
    this.dialogRef.close();
  }
}
