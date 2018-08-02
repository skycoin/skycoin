import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { TranslateService } from '@ngx-translate/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';

import { AppService } from './services/app.service';
import { WalletService } from './services/wallet.service';
import { BigErrorMsgComponent } from './components/layout/big-error-msg/big-error-msg.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(
    private appService: AppService,
    walletService: WalletService,
    dialog: MatDialog,
    translateService: TranslateService,
  ) {
    translateService.setDefaultLang('en');
    translateService.use('en');

    walletService.initialLoadFailed.subscribe(failed => {
      if (failed) {
        const config = new MatDialogConfig();
        config.maxWidth = '100%';
        config.width = '100%';
        config.height = '100%';
        config.disableClose = true;
        config.hasBackdrop = false;

        dialog.open(BigErrorMsgComponent, config);
      }
    });
  }

  ngOnInit() {
    this.appService.testBackend();
  }
}
