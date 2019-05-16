import { TestBed, inject } from '@angular/core/testing';

import { TxEncoder } from './tx-encoder';

describe('TxEncoder', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TxEncoder],
    });
  });

  it('should be created', inject([TxEncoder], (service: TxEncoder) => {
    expect(service).toBeTruthy();
  }));
});
