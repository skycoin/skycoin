import { Subscription, of, Observable, ReplaySubject } from 'rxjs';
import { delay,  map, mergeMap } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';

import { ApiService } from './api.service';
import { BalanceAndOutputsService } from './wallet-operations/balance-and-outputs.service';

/**
 * Basic info of the last block added to the blockchain.
 */
export interface BasicBlockInfo {
  seq: number;
  timestamp: number;
  hash: string;
}

/**
 * Data about the current and max coin supply.
 */
export interface CoinSupply {
  currentSupply: string;
  totalSupply: string;
  currentCoinhourSupply: string;
  totalCoinhourSupply: string;
}

/**
 * Info about the current synchronization state of the blockchain.
 */
export interface ProgressEvent {
  currentBlock: number;
  highestBlock: number;
  synchronized: boolean;
}

/**
 * Allows to check the current state of the blockchain.
 */
@Injectable()
export class BlockchainService {
  private progressSubject: ReplaySubject<ProgressEvent> = new ReplaySubject<ProgressEvent>(1);
  private lastCurrentBlock = 0;
  private lastHighestBlock = 0;
  private nodeSynchronized = false;
  /**
   * Allows the service to update the balance the first time the blockchain state is updated.
   */
  private refreshedBalance = false;

  private dataSubscription: Subscription;

  /**
   * Time interval in which periodic data updates will be made.
   */
  private readonly updatePeriod = 2 * 1000;
  /**
   * Time interval in which the periodic data updates will be restarted after an error.
   */
  private readonly errorUpdatePeriod = 2 * 1000;

  /**
   * allows to know the current synchronization state of the blockchain.
   */
  get progress(): Observable<ProgressEvent> {
    return this.progressSubject.asObservable();
  }

  constructor(
    private apiService: ApiService,
    private ngZone: NgZone,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) {
    this.startDataRefreshSubscription(0);
  }

  /**
   * Gets the basic info of the last block added to the blockchain.
   */
  getLastBlock(): Observable<BasicBlockInfo> {
    return this.apiService.get('last_blocks', { num: 1 }).pipe(map(blocks => {
      return {
        seq: blocks.blocks[0].header.seq,
        timestamp: blocks.blocks[0].header.timestamp,
        hash: blocks.blocks[0].header.block_hash,
      };
    }));
  }

  /**
   * Gets info about the coin supply of the blockchain.
   */
  getCoinSupply(): Observable<CoinSupply> {
    return this.apiService.get('coinSupply').pipe(map(supply => {
      return {
        currentSupply: supply.current_supply,
        totalSupply: supply.total_supply,
        currentCoinhourSupply: supply.current_coinhour_supply,
        totalCoinhourSupply: supply.total_coinhour_supply,
      };
    }));
  }

  /**
   * Makes the service start periodically checking the synchronization state of the blockchain.
   * If this function was called before, the previous procedure is cancelled.
   * @param delayMs Delay before starting to check the data.
   */
  private startDataRefreshSubscription(delayMs: number) {
    if (this.dataSubscription) {
      this.dataSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.dataSubscription = of(0).pipe(
        delay(delayMs),
        mergeMap(() => this.apiService.get('blockchain/progress')),
      ).subscribe((response: any) => {
        this.ngZone.run(() => {
          // Stop if a value is not valid.
          if (!response || !response.current || !response.highest || response.highest === 0 || response.current < this.lastCurrentBlock || response.highest < this.lastHighestBlock) {
            this.startDataRefreshSubscription(this.errorUpdatePeriod);

            return;
          }

          this.lastCurrentBlock = response.current;
          this.lastHighestBlock = response.highest;

          if (response.current === response.highest && !this.nodeSynchronized) {
            this.nodeSynchronized = true;
            this.balanceAndOutputsService.refreshBalance();
            this.refreshedBalance = true;
          } else if (response.current !== response.highest && this.nodeSynchronized) {
            this.nodeSynchronized = false;
          }

          // Refresh the balance the first time the info is retrieved.
          if (!this.refreshedBalance) {
            this.balanceAndOutputsService.refreshBalance();
            this.refreshedBalance = true;
          }

          this.progressSubject.next({
            currentBlock: this.lastCurrentBlock,
            highestBlock: this.lastHighestBlock,
            synchronized: this.nodeSynchronized,
          });

          this.startDataRefreshSubscription(this.updatePeriod);
        });
      }, () => {
        this.startDataRefreshSubscription(this.errorUpdatePeriod);
      });
    });
  }
}
