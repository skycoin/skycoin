import { Component, OnInit } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { LanguageData, LanguageService } from '../../../services/language.service';

/**
 * Allows to change the language displayed by the UI.
 */
@Component({
  selector: 'app-select-language',
  templateUrl: './select-language.component.html',
  styleUrls: ['./select-language.component.scss'],
})
export class SelectLanguageComponent implements OnInit {
  languages: LanguageData[];
  disableDismiss: boolean;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   * @param disableClose Disables the options for closing the modal window without
   * selecting a langhuage.
   */
  public static openDialog(dialog: MatDialog, disableClose = false): MatDialogRef<SelectLanguageComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.disableClose = disableClose;
    config.width = '600px';

    return dialog.open(SelectLanguageComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<SelectLanguageComponent>,
    private languageService: LanguageService,
  ) { }

  ngOnInit() {
    this.disableDismiss = this.dialogRef.disableClose;
    this.languages = this.languageService.languages;
  }

  closePopup(language: LanguageData = null) {
    if (language) {
      this.languageService.changeLanguage(language.code);
    }

    this.dialogRef.close();
  }
}
