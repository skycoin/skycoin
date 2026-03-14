import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, Subscription, of } from 'rxjs';
import { map, mergeMap, delay } from 'rxjs';
import BigNumber from 'bignumber.js';

import { ApiService } from './api.service';
import { BalanceService } from './wallet/balance.service';
import { ConnectionError } from '../enums/connection-error.enum';
import { CoinService } from './coin.service';
import { environment } from '../../environments/environment';
import { GlobalsService } from './globals.service';
import { isEqualOrSuperiorVersion } from '../utils/semver';

export enum ProgressStates {
  Progress,
  Error,
  Restarting,
}

export class ProgressEvent {
  state: ProgressStates;
  error?: ConnectionError;
  currentBlock?: number;
  highestBlock?: number;
}

@Injectable()
export class BlockchainService {
  private progressSubject: BehaviorSubject<ProgressEvent> = new BehaviorSubject<ProgressEvent>(null);
  private synchronizedSubject: BehaviorSubject<any> = new BehaviorSubject<boolean>(false);
  private connectionsSubscription: Subscription;
  private maxDecimals = 6;
  private readonly defaultPeriod = 90000;
  private readonly shortPeriod = 5000;
  private burnRateInternal = new BigNumber(0.5);

  get progress(): Observable<ProgressEvent> {
    return this.progressSubject.asObservable();
  }

  get synchronized() {
    return this.synchronizedSubject.asObservable();
  }

  get currentMaxDecimals(): number {
    return this.maxDecimals;
  }

  get burnRate() {
    return this.burnRateInternal;
  }

  constructor (
    private apiService: ApiService,
    private balanceService: BalanceService,
    private globalsService: GlobalsService,
    coinService: CoinService
  ) {
    coinService.currentCoin.subscribe(() => {
      this.startCheckingNode();
    });
  }

  lastBlock(): Observable<any> {
    return this.apiService.get('last_blocks', { num: 1 }).pipe(
      map(response => response.blocks[0]));
  }

  coinSupply(): Observable<any> {
    return this.apiService.get('coinSupply');
  }

  private startCheckingNode() {
    if (!!this.connectionsSubscription && !this.connectionsSubscription.closed) {
      this.connectionsSubscription.unsubscribe();
    }

    this.balanceService.stopGettingBalances();

    this.globalsService.setNodeVersion(null);
    this.progressSubject.next({ state: ProgressStates.Restarting });
    this.synchronizedSubject.next(false);

    this.connectionsSubscription = this.checkConnectionState()
      .subscribe(
        () => this.checkBlockchainProgress(0),
        error => error && error.reported ? null : this.onLoadBlockchainError()
      );
  }

  private checkBlockchainProgress(delayMs: number) {
    this.connectionsSubscription = of(1).pipe(
      delay(delayMs),
      mergeMap(() => this.getBlockchainProgress()),
    )
      .subscribe(
        (response: any) => this.onBlockchainProgress(response),
        () => this.onLoadBlockchainError()
      );
  }

  private getBlockchainProgress() {
    return this.apiService.get('blockchain/progress');
  }

  private onBlockchainProgress(response: any) {
    this.progressSubject.next({
      state: ProgressStates.Progress,
      currentBlock: response.current,
      highestBlock: response.highest
    });

    if (response.current !== response.highest) {
      this.checkBlockchainProgress(
        response.highest - response.current <= 5 ? this.shortPeriod : this.defaultPeriod
      );
    } else {
      this.synchronizedSubject.next(true);
      this.balanceService.startGettingBalances();
    }
  }

  private onLoadBlockchainError(error: ConnectionError = ConnectionError.UNAVAILABLE_BACKEND) {
    this.progressSubject.next({ state: ProgressStates.Error, error: error });
  }

  private checkConnectionState(): Observable<any> {
    return this.apiService.get('network/connections').pipe(
      map((status: any) => {
        if ((!status.connections || status.connections.length === 0) && !environment.e2eTest) {
          this.onLoadBlockchainError(ConnectionError.NO_ACTIVE_CONNECTIONS);
          throw { reported: true };
        }
      }),
      mergeMap(() => this.apiService.get('health')),
      map((response: any) => {
        this.globalsService.setNodeVersion(response.version.version);
        if (isEqualOrSuperiorVersion(response.version.version, '0.25.0')) {
          this.maxDecimals = response.user_verify_transaction.max_decimals;
          this.burnRateInternal = new BigNumber(response.user_verify_transaction.burn_factor);
        } else {
          this.maxDecimals = 6;
          this.burnRateInternal = new BigNumber(2);
        }
      }),
    );
  }
}
