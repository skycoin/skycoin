import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class PurchaseService {

  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);

  constructor(
    private http: Http,
  ) {
  }

  all() {
    return this.purchaseOrders.asObservable();
  }

  private get(url) {
    return this.http.get(this.purchaseUrl + url)
      .map((res: any) => res.json())
  }

  private post(url, parameters = {}) {
    return this.http.post(this.purchaseUrl + url, parameters)
      .map((res: any) => res.json())
  }
}
