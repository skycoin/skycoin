import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { WalletService } from './wallet.service';
import 'rxjs/add/observable/timer';

@Injectable()
export class BlockchainService {
  private progressSubject: Subject<any> = new BehaviorSubject<any>(null);
  private refreshedBalance = false;

  get progress() {
    return this.progressSubject.asObservable();
  }

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
  ) {
    Observable.timer(0, 2000)
      .flatMap(() => this.getBlockchainProgress())
      .takeWhile((response: any) => !response.current || response.current !== response.highest)
      .subscribe(
        response => {
          this.progressSubject.next(response);

          if (! this.refreshedBalance) {
            this.walletService.refreshBalances();
            this.refreshedBalance = true;
          }
        },
        error => console.log(error),
        () => this.completeLoading(),
      );
  }

  block(id): Observable<any> {
    return this.apiService.get('blocks', { start: id, end: id }).map(response => response.blocks[0]).flatMap(block => {
      return Observable.forkJoin(block.body.txns.map(transaction => {
        if (transaction.inputs && !transaction.inputs.length) {
          return Observable.of(transaction);
        }

        return Observable.forkJoin(transaction.inputs.map(input => this.retrieveInputAddress(input).map(response => {
          return response.owner_address;
        }))).map(inputs => {
          transaction.inputs = inputs;

          return transaction;
        });
      })).map(transactions => {
        block.body.txns = transactions;

        return block;
      });
    });
  }

  blocks(num: number = 5100) {
    return this.apiService.get('last_blocks', { num: num }).map(response => response.blocks.reverse());
  }

  lastBlock() {
    return this.blocks(1).map(blocks => blocks[0]);
  }

  getBlockchainProgress() {
    return this.apiService.get('blockchain/progress');
  }

  coinSupply() {
    return this.apiService.get('coinSupply');
  }

  private completeLoading() {
    this.progressSubject.next({ current: 999999999999, highest: 999999999999 });
    this.walletService.refreshBalances();
  }

  private retrieveInputAddress(input: string) {
    return this.apiService.get('uxout', {uxid: input});
  }
}
