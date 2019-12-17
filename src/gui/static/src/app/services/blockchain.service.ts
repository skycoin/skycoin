import { timer as observableTimer, throwError as observableThrowError, BehaviorSubject } from 'rxjs';
import { retryWhen, concat, delay, exhaustMap, take, map } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { WalletService } from './wallet.service';

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
    this.apiService.get('health').pipe(retryWhen(errors => errors.pipe(delay(1000), take(10), concat(observableThrowError('')))))
      .subscribe ((response: any) => this.maxDecimals = response.user_verify_transaction.max_decimals);

    this.ngZone.runOutsideAngular(() => {
      observableTimer(0, 2000).pipe(
        exhaustMap(() => this.getBlockchainProgress()),
        retryWhen(errors => errors.pipe(delay(2000))))
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
        );
    });
  }

  lastBlock() {
    return this.apiService.get('last_blocks', { num: 1 }).pipe(map(blocks => blocks.blocks[0]));
  }

  getBlockchainProgress() {
    return this.apiService.get('blockchain/progress');
  }

  coinSupply() {
    return this.apiService.get('coinSupply');
  }
}
