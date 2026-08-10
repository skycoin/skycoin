import { Component, HostListener, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';

import { CoinService } from '../../../services/coin.service';
import { BaseCoin } from '../../../coins/basecoin';
import { SpendingService } from '../../../services/wallet/spending.service';
import { MsgBarService } from '../../../services/msg-bar.service';

@Component({
    selector: 'app-select-coin-overlay',
    templateUrl: './select-coin-overlay.component.html',
    styleUrls: ['./select-coin-overlay.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class SelectCoinOverlayComponent implements OnDestroy {

  coins: BaseCoin[] = [];

  constructor(
    private dialogRef: MatDialogRef<SelectCoinOverlayComponent>,
    private coinService: CoinService,
    private spendingService: SpendingService,
    private translate: TranslateService,
    private msgBarService: MsgBarService,
  ) {
    this.coins = this.coinService.coins;
  }

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent) {
    if (event.keyCode === 27) {
      this.close(null);
    }
  }

  close(result: BaseCoin | null) {
    this.msgBarService.hide();
    if (result && this.spendingService.isInjectingTx) {
      this.msgBarService.showError('change-coin.injecting-tx');
      return;
    }
    this.dialogRef.close(result);
  }
}
