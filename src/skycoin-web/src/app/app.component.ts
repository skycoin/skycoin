import { Component, OnInit, Renderer2, ViewChild } from '@angular/core';
import { LanguageService } from './services/language.service';
import { Router, RouterEvent, NavigationEnd } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';

import { config } from './app.config';
import { environment } from '../environments/environment';
import { CipherProvider, InitializationResults } from './services/cipher.provider';
import { CustomMatDialogService } from './services/custom-mat-dialog.service';
import { Bip39WordListService } from './services/bip39-word-list.service';
import { MsgBarComponent } from './components/layout/msg-bar/msg-bar.component';
import { MsgBarService } from './services/msg-bar.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  @ViewChild('msgBar') msgBar: MsgBarComponent;

  current: number;
  highest: number;
  otcEnabled: boolean;
  version: string;
  browserCompatibleWithWasm = true;
  wasmFileLoaded = true;

  constructor(
    private languageService: LanguageService,
    cipherProvider: CipherProvider,
    router: Router,
    dialog: CustomMatDialogService,
    renderer: Renderer2,
    private bip38WordList: Bip39WordListService,
    private msgBarService: MsgBarService,
  ) {
    router.events.subscribe((event: RouterEvent) => {
      if (event instanceof NavigationEnd) {
        window.scrollTo(0, 0);
      }
    });

    cipherProvider.initialize().subscribe(response => {
      this.checkCipherProviderResponse(response);
    }, response => this.checkCipherProviderResponse(response));

    dialog.showingDialog.subscribe(value => {
      if (!value) {
        renderer.addClass(document.body, 'fix-error-position');
      } else {
        renderer.removeClass(document.body, 'fix-error-position');
      }
    });
  }

  ngOnInit() {
    this.otcEnabled = config.otcEnabled;
    this.languageService.loadLanguageSettings();

    window.onbeforeunload = (e) => {
      if (environment.production && !environment.e2eTest && !window['isElectron']) {
        e.preventDefault();
        e.returnValue = '';
      }
    };

    this.msgBarService.msgBarComponent = this.msgBar;
  }

  loading() {
    return !this.current || !this.highest || this.current !== this.highest;
  }

  private checkCipherProviderResponse(response) {
    if (window['removeSplash']) {
      setTimeout(() => window['removeSplash']());
    }
    if (response !== InitializationResults.Ok) {
      if (response === InitializationResults.ErrorLoadingWasmFile) {
        this.wasmFileLoaded = false;
      } else {
        this.browserCompatibleWithWasm = false;
      }
    }
  }
}
