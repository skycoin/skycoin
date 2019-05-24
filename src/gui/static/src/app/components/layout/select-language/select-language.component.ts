import { Component, OnInit } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

import { LanguageData, LanguageService } from '../../../services/language.service';

@Component({
  selector: 'app-select-language',
  templateUrl: './select-language.component.html',
  styleUrls: ['./select-language.component.scss'],
})
export class SelectLanguageComponent implements OnInit {

  languages: LanguageData[];
  disableDismiss: boolean;

  constructor(
    public dialogRef: MatDialogRef<SelectLanguageComponent>,
    private languageService: LanguageService,
  ) { }

  ngOnInit() {
    this.disableDismiss = this.dialogRef.disableClose;
    this.languages = this.languageService.languages;
  }

  closePopup(language: LanguageData = null) {
    this.dialogRef.close(language ? language.code : undefined);
  }
}
