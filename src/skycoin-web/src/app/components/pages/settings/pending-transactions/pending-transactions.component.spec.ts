import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { PendingTransactionsComponent } from './pending-transactions.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { HistoryService } from '../../../../services/wallet/history.service';
import { NavBarService } from '../../../../services/nav-bar.service';
import { CoinService } from '../../../../services/coin.service';
import { GlobalsService } from '../../../../services/globals.service';
import {
  MockTranslatePipe,
  MockDateTimePipe,
  MockNavBarService,
  MockCoinService,
  MockWalletService,
  MockHistoryService,
  MockGlobalsService
} from '../../../../utils/test-mocks';

describe('PendingTransactionsComponent', () => {
  let component: PendingTransactionsComponent;
  let fixture: ComponentFixture<PendingTransactionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        PendingTransactionsComponent,
        MockDateTimePipe,
        MockTranslatePipe
      ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: WalletService, useClass: MockWalletService },
        { provide: HistoryService, useClass: MockHistoryService },
        { provide: NavBarService, useClass: MockNavBarService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService }
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PendingTransactionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
