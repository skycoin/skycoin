/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { TransactionDetailService } from './transaction-detail.service';

describe('TransactionDetailService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TransactionDetailService]
    });
  });

  it('should ...', inject([TransactionDetailService], (service: TransactionDetailService) => {
    expect(service).toBeTruthy();
  }));
});
