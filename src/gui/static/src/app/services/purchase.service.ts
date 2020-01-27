import { Injectable } from '@angular/core';
import { Subject, BehaviorSubject, Observable } from 'rxjs';
import { PurchaseOrder, TellerConfig, Wallet } from '../app.datatypes';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { map, mergeMap } from 'rxjs/operators';
import { WalletBase } from './wallet-operations/wallet-objects';
import { WalletsAndAddressesService } from './wallet-operations/wallets-and-addresses.service';

@Injectable()
export class PurchaseService {
  private configSubject: Subject<TellerConfig> = new BehaviorSubject<TellerConfig>(null);
  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);
  private purchaseUrl = environment.tellerUrl;

  constructor(
    private httpClient: HttpClient,
    private walletsAndAddressesService: WalletsAndAddressesService,
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
    return this.get('config').pipe(
      map((response: any) => ({
        enabled: true,
        sky_btc_exchange_rate: parseFloat(response.sky_btc_exchange_rate),
      })))
      .subscribe(response => this.configSubject.next(response));
  }

  generate(wallet: WalletBase): Observable<PurchaseOrder> {
    return this.walletsAndAddressesService.addAddressesToWallet(wallet, 1).pipe(mergeMap(address => {
      return this.post('bind', { skyaddr: address[0].address, coin_type: 'BTC' }).pipe(
        map(response => ({
          coin_type: response.coin_type,
          deposit_address: response.deposit_address,
          filename: wallet.id,
          recipient_address: address[0].address,
          status: 'waiting_deposit',
        })));
    }));
  }

  scan(address: string) {
    return this.get('status?skyaddr=' + address).pipe(
      map((response: any) => {
        if (!response.statuses || response.statuses.length > 1) {
          throw new Error('too many purchase orders found');
        }

        return response.statuses[0];
      }));
  }

  private get(url): Observable<any> {
    return this.httpClient.get(this.purchaseUrl + url);
  }

  private post(url, parameters = {}): Observable<any> {
    return this.httpClient.post(this.purchaseUrl + url, parameters);
  }
}
