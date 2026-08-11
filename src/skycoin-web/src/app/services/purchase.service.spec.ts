import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { PurchaseService } from './purchase.service';

describe('PurchaseService', () => {
  let service: PurchaseService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PurchaseService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    service = TestBed.inject(PurchaseService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
