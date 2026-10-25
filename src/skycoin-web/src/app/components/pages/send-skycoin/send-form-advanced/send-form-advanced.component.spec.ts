import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { UntypedFormBuilder } from '@angular/forms';
import { of } from 'rxjs';

import { SendFormAdvancedComponent } from './send-form-advanced.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { CoinService } from '../../../../services/coin.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { NavBarService } from '../../../../services/nav-bar.service';
import { PriceService } from '../../../../services/price.service';
import { MockTranslatePipe,
  MockWalletService,
  MockSpendingService,
  MockCoinService,
  MockBlockchainService,
  MockCustomMatDialogService,
  MockNavBarService,
  MockPriceService,
  MockMsgBarService} from '../../../../utils/test-mocks';
import { MsgBarService } from '../../../../services/msg-bar.service';

describe('SendFormAdvancedComponent', () => {
  let component: SendFormAdvancedComponent;
  let fixture: ComponentFixture<SendFormAdvancedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ SendFormAdvancedComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        UntypedFormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: SpendingService, useClass: MockSpendingService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        { provide: NavBarService, useClass: MockNavBarService },
        { provide: PriceService, useClass: MockPriceService },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendFormAdvancedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  function bulk(result: any) {
    const dialog = TestBed.inject(CustomMatDialogService) as any;
    dialog.open = jasmine.createSpy('open').and.returnValue({ afterClosed: () => of(result) });
    component.openBulkSend();

    return dialog.open.calls.mostRecent().args[1].data;
  }

  it('fills the destinations from bulk text', () => {
    bulk([{ address: 'a1', coins: '1' }, { address: 'a2', coins: '2' }]);
    expect(component.destControls.map(c => c.value)).toEqual([
      { address: 'a1', coins: '1', hours: '' },
      { address: 'a2', coins: '2', hours: '' },
    ]);
    expect(component.autoHours).toBeTrue();
  });

  it('turns off automatic hours when the bulk text has hours', () => {
    bulk([{ address: 'a1', coins: '1', hours: '3' }]);
    expect(component.destControls[0].value).toEqual({ address: 'a1', coins: '1', hours: '3' });
    expect(component.autoHours).toBeFalse();
  });

  it('offers the current destinations as text and keeps them on cancel', () => {
    bulk([{ address: 'a1', coins: '1' }, { address: 'a2', coins: '2' }]);
    expect(bulk(undefined)).toBe('a1,1\na2,2');
    expect(component.destControls.length).toBe(2);
  });
});
