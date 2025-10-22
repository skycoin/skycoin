import { TestBed, inject } from '@angular/core/testing';
import { TranslateService } from '@ngx-translate/core';

import { CoinService } from './coin.service';
import { MockTranslateService } from '../utils/test-mocks';

describe('CoinService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        CoinService,
        { provide: TranslateService, useClass: MockTranslateService },
      ]
    });
  });

  it('should be created', inject([CoinService], (service: CoinService) => {
    expect(service).toBeTruthy();
  }));
});
