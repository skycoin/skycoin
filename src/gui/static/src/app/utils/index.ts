import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material';

import { ConfirmationData } from '../app.datatypes';
import { ConfirmationComponent } from '../components/layout/confirmation/confirmation.component';

export function showConfirmationModal(dialog: MatDialog, confirmationData: ConfirmationData): MatDialogRef<ConfirmationComponent, any> {
  return dialog.open(ConfirmationComponent, <MatDialogConfig> {
    width: '450px',
    data: confirmationData,
    autoFocus: false,
  });
}
