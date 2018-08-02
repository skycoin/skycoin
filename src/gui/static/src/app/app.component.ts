import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { TranslateService } from '@ngx-translate/core';

import { AppService } from './services/app.service';
import { WalletService } from './services/wallet.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(
    private appService: AppService,
    walletService: WalletService,
    translateService: TranslateService,
  ) {
    translateService.setDefaultLang('en');
    translateService.use('en');

    walletService.initialLoadFailed.subscribe(failed => {
      if (failed) {
        // The "?2" part indicates that error number 2 should be displayed.
        window.location.assign('assets/error-alert/index.html?2');
      }
    });
  }

  ngOnInit() {
    this.appService.testBackend();
  }
}
