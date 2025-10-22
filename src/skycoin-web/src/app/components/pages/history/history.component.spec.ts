import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { FormBuilder } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Observable } from 'rxjs';

import { HistoryComponent } from './history.component';
import { HistoryService } from '../../../services/wallet/history.service';
import { PriceService } from '../../../services/price.service';
import { CoinService } from '../../../services/coin.service';
import { MockTranslatePipe, MockPriceService, MockCoinService, MockHistoryService, MockCustomMatDialogService, MockDateTimePipe, MockWalletService } from '../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { WalletService } from '../../../services/wallet/wallet.service';

describe('HistoryComponent', () => {
  let component: HistoryComponent;
  let fixture: ComponentFixture<HistoryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HistoryComponent, MockTranslatePipe, MockDateTimePipe ],
      providers: [
        FormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: ActivatedRoute, useValue: { queryParams: Observable.of({}) } },
        { provide: HistoryService, useClass: MockHistoryService },
        { provide: PriceService, useClass: MockPriceService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ],
      schemas: [ NO_ERRORS_SCHEMA ]
    });
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HistoryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
