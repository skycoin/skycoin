import { Component, Input, OnInit } from '@angular/core';
import { ExchangeOrder } from '../../../../app.datatypes';
import { ExchangeService } from '../../../../services/exchange.service';
import { TranslateService } from '@ngx-translate/core';
import { QrCodeComponent } from '../../../layout/qr-code/qr-code.component';
import { MatDialog, MatDialogConfig } from '@angular/material';

@Component({
  selector: 'app-exchange-status',
  templateUrl: './exchange-status.component.html',
  styleUrls: ['./exchange-status.component.scss'],
})
export class ExchangeStatusComponent implements OnInit {
  @Input() order: ExchangeOrder;

  readonly statuses = [
    'user_waiting',
    'market_waiting_confirmations',
    'market_confirmed',
    'market_exchanged',
    'market_withdraw_waiting',
    'complete',
    'error',
  ];

  get fromCoin() {
    return this.order.pair.split('/')[0].toUpperCase();
  }

  get toCoin() {
    return this.order.pair.split('/')[1].toUpperCase();
  }

  get translatedStatus() {
    const status = this.order.status.replace(/_/g, '-');
    const params = {
      from: this.fromCoin,
      amount: this.order.fromAmount,
      to: this.toCoin,
    };

    return this.translateService.instant(`exchange.statuses.${status}`, params);
  }

  get statusIcon() {
    if (this.order.status === this.statuses[5]) {
      return 'done';
    }

    if (this.order.status === this.statuses[6]) {
      return 'close';
    }

    return 'refresh';
  }

  get progress() {
    let index = this.statuses.indexOf(this.order.status);

    index = Math.min(index, 5) + 1;

    return Math.ceil((100 / 6) * index);
  }

  constructor(
    private exchangeService: ExchangeService,
    private translateService: TranslateService,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    const fromAmount = this.order.fromAmount;

    this.exchangeService.status(this.order.id).subscribe(order => {
      this.order = { ...order, fromAmount };
    });
  }

  showQrCode(event, address) {
    event.stopPropagation();

    const config = new MatDialogConfig();
    config.data = { address };
    this.dialog.open(QrCodeComponent, config);
  }
}
