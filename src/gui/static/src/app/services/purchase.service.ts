import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class PurchaseService {

  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);

  all() {
    return this.purchaseOrders.asObservable();
  }
}
