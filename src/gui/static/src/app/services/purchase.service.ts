import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { PurchaseOrder, TellerConfig, Wallet } from '../app.datatypes';
import { WalletService } from './wallet.service';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class PurchaseService {
  private configSubject: Subject<TellerConfig> = new BehaviorSubject<TellerConfig>(null);
  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);
  private purchaseUrl = environment.tellerUrl;

  constructor(
    private httpClient: HttpClient,
    private walletService: WalletService,
  ) {
    this.getConfig();
  }

  all() {
    return this.purchaseOrders.asObservable();
  }

  config(): Observable<TellerConfig> {
    return this.configSubject.asObservable();
  }

  getConfig() {
    return this.get('config')
      .map((response: any) => ({
        enabled: true,
        sky_btc_exchange_rate: parseFloat(response.sky_btc_exchange_rate),
      }))
      .subscribe(response => this.configSubject.next(response));
  }

  generate(wallet: Wallet): Observable<PurchaseOrder> {
    return this.walletService.addAddress(wallet).flatMap(address => {
      return this.post('bind', { skyaddr: address.address, coin_type: 'BTC' })
        .map(response => ({
          coin_type: response.coin_type,
          deposit_address: response.deposit_address,
          filename: wallet.filename,
          recipient_address: address.address,
          status: 'waiting_deposit',
        }));
    });
  }

  scan(address: string) {
    return this.get('status?skyaddr=' + address)
      .map((response: any) => {
        if (!response.statuses || response.statuses.length > 1) {
          throw new Error('too many purchase orders found');
        }

        return response.statuses[0];
      });
  }

  private get(url): any {
    return this.httpClient.get(this.purchaseUrl + url);
  }

  private post(url, parameters = {}): any {
    return this.httpClient.post(this.purchaseUrl + url, parameters);
  }
}
