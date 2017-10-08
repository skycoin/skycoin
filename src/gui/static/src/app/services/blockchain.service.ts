import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class BlockchainService {

  constructor(
    private apiService: ApiService,
  ) { }

  addressTransactions(id): Observable<any> {
    return this.apiService.get('explorer/address', { address: id });
  }

  addressBalance(id): Observable<any> {
    return this.apiService.get('outputs', { addrs: id });
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

  progress() {
    return this.apiService.get('blockchain/progress');
  }

  private retrieveInputAddress(input: string) {
    return this.apiService.get('uxout', {uxid: input});
  }
}
