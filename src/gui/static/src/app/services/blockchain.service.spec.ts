import { TestBed, inject } from '@angular/core/testing';

import { BlockchainService } from './blockchain.service';

describe('BlockchainService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BlockchainService]
    });
  });

  it('should be created', inject([BlockchainService], (service: BlockchainService) => {
    expect(service).toBeTruthy();
  }));
});
