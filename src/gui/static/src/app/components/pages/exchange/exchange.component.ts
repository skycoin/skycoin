import { Component, OnInit } from '@angular/core';
import { ExchangeOrder } from '../../../app.datatypes';
import { ExchangeService } from '../../../services/exchange.service';

@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent implements OnInit {
  order: ExchangeOrder;

  constructor(
    private exchangeService: ExchangeService,
  ) { }

  ngOnInit() {
    const lastOrder = this.exchangeService.lastOrder;

    if (lastOrder) {
      if (this.exchangeService.isOrderFinished(lastOrder)) {
        this.showLast();
      }
    }
  }

  showLast() {
    this.order = this.exchangeService.lastOrder;
  }

  showStatus(lastOrder) {
    this.order = lastOrder;
  }
}
