import { TestBed } from '@angular/core/testing';
import { Observable, of, throwError } from 'rxjs';

import { NodeNotices, NodeNoticesService } from './node-notices.service';
import { ApiService } from './api.service';
import { CoinService } from './coin.service';
import { MockCoinService } from '../utils/test-mocks';

describe('NodeNoticesService', () => {
  let health: Observable<any>;

  function notices(): NodeNotices | undefined {
    TestBed.configureTestingModule({
      providers: [
        NodeNoticesService,
        { provide: CoinService, useClass: MockCoinService },
        { provide: ApiService, useValue: { get: () => health } },
      ],
    });
    let result: NodeNotices | undefined;
    TestBed.inject(NodeNoticesService).notices.subscribe(n => result = n);

    return result;
  }

  it('warns when the node has CSRF turned off', () => {
    health = of({ csrf_enabled: false });
    expect(notices()).toEqual({ csrfDisabled: true });
  });

  it('says nothing when the node checks CSRF', () => {
    health = of({ csrf_enabled: true });
    expect(notices()).toEqual({ csrfDisabled: false });
  });

  it('fails quietly when the node cannot be reached', () => {
    health = throwError(() => new Error('unreachable'));
    expect(notices()).toEqual({ csrfDisabled: false });
  });
});
