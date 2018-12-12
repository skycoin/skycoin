import { Component } from '@angular/core';
import { ExchangeOrder } from '../../../app.datatypes';
import { ExchangeService } from '../../../services/exchange.service';

@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent {
  order: ExchangeOrder;

  constructor(
    private exchangeService: ExchangeService,
  ) { }

  showStatus(lastOrder) {
    this.order = lastOrder;
  }

  hasLast() {
    return !!this.exchangeService.getLastOrder();
  }

  showLast(event) {
    event.preventDefault();

    this.order = this.exchangeService.getLastOrder();
  }
}
