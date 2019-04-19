import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { WalletService } from './wallet.service';
import 'rxjs/add/observable/timer';
import 'rxjs/add/operator/retryWhen';
import 'rxjs/add/operator/concat';
import 'rxjs/add/operator/exhaustMap';

@Injectable()
export class BlockchainService {
  private progressSubject: BehaviorSubject<any> = new BehaviorSubject<any>(null);
  private synchronizedSubject: BehaviorSubject<any> = new BehaviorSubject<boolean>(false);
  private refreshedBalance = false;
  private lastCurrentBlock = 0;
  private lastHighestBlock = 0;
  private maxDecimals = 6;

  get progress() {
    return this.progressSubject.asObservable();
  }

  get currentMaxDecimals(): number {
    return this.maxDecimals;
  }

  get synchronized() {
    return this.synchronizedSubject.asObservable();
  }

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
    private ngZone: NgZone,
  ) {
    this.apiService.get('health').retryWhen(errors => errors.delay(1000).take(10).concat(Observable.throw('')))
      .subscribe ((response: any) => this.maxDecimals = response.user_verify_transaction.max_decimals);

    this.ngZone.runOutsideAngular(() => {
      Observable.timer(0, 2000)
        .exhaustMap(() => this.getBlockchainProgress())
        .retryWhen(errors => errors.delay(2000))
        .subscribe(
          response => this.ngZone.run(() => {
            if (!response.current || !response.highest || response.current < this.lastCurrentBlock || response.highest < this.lastHighestBlock) {
              return;
            }

            this.lastCurrentBlock = response.current;
            this.lastHighestBlock = response.highest;

            this.progressSubject.next(response);

            if (!this.refreshedBalance) {
              this.walletService.refreshBalances();
              this.refreshedBalance = true;
            }

            if (response.current === response.highest && !this.synchronizedSubject.value) {
              this.synchronizedSubject.next(true);
              this.walletService.refreshBalances();
            } else if (response.current !== response.highest && this.synchronizedSubject.value) {
              this.synchronizedSubject.next(false);
            }
          }),
          error => console.log(error),
        );
    });
  }

  lastBlock() {
    return this.apiService.get('last_blocks', { num: 1 }).map(blocks => blocks.blocks[0]);
  }

  getBlockchainProgress() {
    return this.apiService.get('blockchain/progress');
  }

  coinSupply() {
    return this.apiService.get('coinSupply');
  }
}
