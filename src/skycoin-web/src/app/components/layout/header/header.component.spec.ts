import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { HeaderComponent } from './header.component';
import { PriceService } from '../../../services/price.service';
import { BalanceService } from '../../../services/wallet/balance.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { CoinService } from '../../../services/coin.service';
import { NodeNotices, NodeNoticesService } from '../../../services/node-notices.service';
import { BehaviorSubject } from 'rxjs';
import {
  MockTranslatePipe,
  MockCoinService,
  MockBlockchainService,
  MockPriceService,
  MockBalanceService
} from '../../../utils/test-mocks';

describe('HeaderComponent', () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;
  let notices: BehaviorSubject<NodeNotices>;

  beforeEach(waitForAsync(() => {
    notices = new BehaviorSubject<NodeNotices>({ csrfDisabled: false });
    TestBed.configureTestingModule({
      declarations: [ HeaderComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: PriceService, useClass: MockPriceService },
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: NodeNoticesService, useValue: { notices } },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('shows no notice while the node has none', () => {
    expect(fixture.nativeElement.querySelector('.-csrf-notice')).toBeNull();
  });

  it('warns when the node has CSRF turned off', () => {
    notices.next({ csrfDisabled: true });
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.-csrf-notice')).not.toBeNull();
  });
});
