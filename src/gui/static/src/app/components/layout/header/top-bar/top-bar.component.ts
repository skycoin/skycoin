import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { AppService } from '../../../../services/app.service';
import { LanguageData, LanguageService } from '../../../../services/language.service';
import { openChangeLanguageModal } from '../../../../utils';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-top-bar',
  templateUrl: './top-bar.component.html',
  styleUrls: ['./top-bar.component.scss'],
})
export class TopBarComponent implements OnInit, OnDestroy {
  @Input() headline: string;

  language: LanguageData;

  private subscription: Subscription;

  constructor(
    public appService: AppService,
    private languageService: LanguageService,
    private dialog: MatDialog,
  ) {}

  ngOnInit() {
    this.subscription = this.languageService.currentLanguage
      .subscribe(lang => this.language = lang);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  changelanguage() {
    openChangeLanguageModal(this.dialog)
      .subscribe(response => {
        if (response) {
          this.languageService.changeLanguage(response);
        }
      });
  }
}
