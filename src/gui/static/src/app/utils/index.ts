import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material';
import { Observable } from 'rxjs/Observable';

import { ConfirmationData } from '../app.datatypes';
import { ConfirmationComponent } from '../components/layout/confirmation/confirmation.component';
import { SelectLanguageComponent } from '../components/layout/select-language/select-language.component';

export function showConfirmationModal(dialog: MatDialog, confirmationData: ConfirmationData): MatDialogRef<ConfirmationComponent, any> {
  return dialog.open(ConfirmationComponent, <MatDialogConfig> {
    width: '450px',
    data: confirmationData,
    autoFocus: false,
  });
}

export function openChangeLanguageModal(dialog: MatDialog, disableClose = false): Observable<any> {
  const config = new MatDialogConfig();
  config.width = '600px';
  config.disableClose = disableClose;
  config.autoFocus = false;

  return dialog.open(SelectLanguageComponent, config).afterClosed();
}
