import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { of } from 'rxjs';

import { HwWalletDaemonService } from './hw-wallet-daemon.service';
import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { HWOperationResults } from '../utils/operation-error';

describe('HwWalletDaemonService word requests', () => {
  const base = 'http://127.0.0.1:9510/api/v1';
  let service: HwWalletDaemonService;
  let http: HttpTestingController;
  let requestWord: jasmine.Spy;

  beforeEach(() => {
    requestWord = jasmine.createSpy('requestWord');
    TestBed.configureTestingModule({
      providers: [
        HwWalletDaemonService,
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: HwWalletPinService, useValue: {} },
        { provide: HwWalletSeedWordService, useValue: { requestWord } },
      ],
    });
    service = TestBed.inject(HwWalletDaemonService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('sends each word the device asks for until it finishes', () => {
    requestWord.and.returnValues(of('abandon'), of('zoo'));
    let result: any;
    service.post('/recovery', { word_count: 12 }).subscribe(r => result = r);

    http.expectOne(base + '/recovery').flush({ data: 'WordRequest' });
    const first = http.expectOne(base + '/intermediate/word');
    expect(first.request.body).toBe(JSON.stringify({ word: 'abandon' }));
    first.flush({ data: 'WordRequest' });
    const second = http.expectOne(base + '/intermediate/word');
    expect(second.request.body).toBe(JSON.stringify({ word: 'zoo' }));
    second.flush({ data: 'Device recovered' });

    expect(result.data).toBe('Device recovered');
  });

  it('cancels the operation when no word is entered', () => {
    requestWord.and.returnValue(of(undefined));
    let error: any;
    service.post('/recovery', {}).subscribe({ error: e => error = e });

    http.expectOne(base + '/recovery').flush({ data: 'WordRequest' });
    const cancel = http.expectOne(base + '/cancel');
    expect(cancel.request.method).toBe('PUT');
    cancel.flush({ data: 'Action cancelled by user' });

    expect(error.type).toBe(HWOperationResults.FailedOrRefused);
  });
});
