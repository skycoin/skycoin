import { TestBed, inject } from '@angular/core/testing';
import { MockBackend } from '@angular/http/testing';
import { XHRBackend } from '@angular/http';
import { TranslateService } from '@ngx-translate/core';
import { HttpClientModule } from '@angular/common/http';

import { ApiService } from './api.service';
import { MockTranslateService, MockCoinService } from '../utils/test-mocks';
import { CoinService } from './coin.service';

describe('ApiService', () => {
  let service: ApiService;
  let mockBackEnd: MockBackend;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientModule ],
      providers: [
        ApiService,
        { provide: XHRBackend, useClass: MockBackend },
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: CoinService, useClass: MockCoinService }
      ]
    });
  });

  beforeEach(inject([ApiService, XHRBackend], (serv, mock) => {
    service = serv;
    mockBackEnd = mock;
  }));

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
