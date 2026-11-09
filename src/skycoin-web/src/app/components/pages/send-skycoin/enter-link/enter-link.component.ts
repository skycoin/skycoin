import { ChangeDetectionStrategy, Component } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs';

import { PaymentLink, parsePaymentLink } from '../../../../utils/qr';

/**
 * Takes a pasted payment link such as skycoin:ADDR?amount=1&hours=2. Closes
 * with the link's data, or with undefined when cancelled. Use EnterLinkComponent.ask.
 */
@Component({
    selector: 'app-enter-link',
    templateUrl: './enter-link.component.html',
    styleUrls: ['./enter-link.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class EnterLinkComponent {
  form: UntypedFormGroup;
  error = '';

  static ask(dialog: MatDialog): Observable<PaymentLink | undefined> {
    return dialog.open(EnterLinkComponent, { width: '566px' }).afterClosed();
  }

  constructor(
    public dialogRef: MatDialogRef<EnterLinkComponent>,
    formBuilder: UntypedFormBuilder,
  ) {
    this.form = formBuilder.group({ link: [''] });
  }

  process() {
    const link = parsePaymentLink(this.form.value.link);
    if (!link) {
      this.error = 'send.fill-with-link.invalid-link-error';

      return;
    }
    this.dialogRef.close(link);
  }

  cancel() {
    this.dialogRef.close();
  }
}
