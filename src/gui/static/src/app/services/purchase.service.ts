import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Http } from '@angular/http';

@Injectable()
export class PurchaseService {

  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);
  // private purchaseUrl: string = 'https://event.skycoin.net/api/';
  private purchaseUrl: string = 'http://127.0.01:7071/api/';
  // private purchaseUrl: string = '/teller/';

  constructor(
    private http: Http,
  ) {
    this.retrievePurchaseOrders();
  }

  all() {
    return this.purchaseOrders.asObservable();
  }

  generate(address: string) {
    return this.post('bind', { skyaddr: address })
      .do(response => {
        this.purchaseOrders.first().subscribe(orders => {
          let index = orders.findIndex(order => order.address === address);
          if (index === -1) {
            orders.push({address: address, addresses: []});
            index = orders.length - 1;
          }
          const timestamp = Math.floor(Date.now() / 1000);
          orders[index].addresses.unshift({
            btc: response.btc_address,
            status: 'waiting_deposit',
            created: timestamp,
            updated: timestamp,
          });
          this.updatePurchaseOrders(orders)
        });
      });
  }

  scan(address: string) {
    return this.get('status?skyaddr=' + address).do(response => {
      this.purchaseOrders.first().subscribe(orders => {
        let index = orders.findIndex(order => order.address === address);
        // Sort addresses ascending by creation date to match teller status response
        orders[index].addresses.sort((a, b) =>  b.created - a.created);
        for (const btcAddress of orders[index].addresses) {
          // Splice last status to assign this to the latest known order
          const status = response.statuses.splice(-1,1)[0];
          btcAddress.status = status.status;
          btcAddress.updated = status.update_at;
        }

        this.updatePurchaseOrders(orders)
      });
    });
  }

  private get(url) {
    return this.http.get(this.purchaseUrl + url)
      .map((res: any) => res.json())
  }

  private post(url, parameters = {}) {
    return this.http.post(this.purchaseUrl + url, parameters)
      .map((res: any) => res.json())
  }

  private retrievePurchaseOrders() {
    const orders = JSON.parse(window.localStorage.getItem('purchaseOrders'));
    if (orders) {
      this.purchaseOrders.next(orders);
    }
  }

  private updatePurchaseOrders(collection: any[]) {
    this.purchaseOrders.next(collection);
    window.localStorage.setItem('purchaseOrders', JSON.stringify(collection));
  }
}
