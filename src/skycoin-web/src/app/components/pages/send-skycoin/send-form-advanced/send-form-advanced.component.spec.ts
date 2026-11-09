import { waitForAsync, ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { BigNumber } from 'bignumber.js';
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
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';

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

  it('fills a destination with the available coins less the others', () => {
    bulk([{ address: 'a1', coins: '2' }, { address: 'a2', coins: '' }]);
    component.availableCoins = new BigNumber('10.5');
    component.assignAll(1);
    expect(component.destControls[1].get('coins')!.value).toBe('8.5');
  });

  it('shows an error when nothing is left for a destination', () => {
    const msgBar = TestBed.inject(MsgBarService);
    spyOn(msgBar, 'showError');
    bulk([{ address: 'a1', coins: '5' }, { address: 'a2', coins: '' }]);
    component.availableCoins = new BigNumber('5');
    component.assignAll(1);
    expect(msgBar.showError).toHaveBeenCalledWith('send.no-coins-left-error');
    expect(component.destControls[1].get('coins')!.value).toBe('');
  });

  describe('unsigned transaction from typed addresses', () => {
    let dialog: any;

    function confirmUnsigned(confirmed: boolean) {
      dialog = TestBed.inject(CustomMatDialogService) as any;
      dialog.open = jasmine.createSpy('open').and.returnValue({ afterClosed: () => of(confirmed) });
      component.changeFormType(DoubleButtonActive.RightButton);
    }

    it('asks before switching and stays on the signed form when cancelled', () => {
      confirmUnsigned(false);
      expect(dialog.open.calls.mostRecent().args[1].data.text).toBe('send.unsigned-confirmation');
      expect(component.manualUnsigned).toBeFalse();
      expect(component.modeButton).toBe(DoubleButtonActive.LeftButton);
    });

    it('needs no wallet and builds the transaction from the typed addresses', fakeAsync(() => {
      const spending = TestBed.inject(SpendingService) as any;
      spending.getAddressesUnspentOutputs = jasmine.createSpy('outputs').and.returnValue(of([
        { address: 'a1', hash: 'h1', coins: new BigNumber(3), calculated_hours: new BigNumber(10) },
        { address: 'a2', hash: 'h2', coins: new BigNumber(2), calculated_hours: new BigNumber(10) },
      ]));
      spending.createTransaction = jasmine.createSpy('create').and.returnValue(of({ encoded: 'unsigned-tx' }));

      confirmUnsigned(true);
      expect(component.manualUnsigned).toBeTrue();
      expect(component.form.get('wallet')!.valid).toBeTrue();

      component.form.get('manualAddresses')!.setValue(' a1 , a2,');
      tick(500);
      expect(spending.getAddressesUnspentOutputs).toHaveBeenCalledWith(['a1', 'a2']);
      expect(component.availableCoins.toString()).toBe('5');

      component.destControls[0].get('address')!.setValue('dest');
      component.destControls[0].get('coins')!.setValue('4');
      expect(component.form.valid).toBeTrue();

      component.unsignedButton = { isLoading: () => false, setLoading: () => undefined, resetState: () => undefined } as any;
      component.createUnsigned();
      const args = spending.createTransaction.calls.mostRecent().args;
      expect(args[0]).toBeNull();
      expect(args[1]).toEqual(['a1', 'a2']);
      expect(args[7]).toBeTrue();
      expect(dialog.open.calls.mostRecent().args[1].data.rawTx).toBe('unsigned-tx');
    }));
  });

  it('fills the first destination from a payment link and sets its hours', () => {
    const dialog = TestBed.inject(CustomMatDialogService) as any;
    dialog.open = jasmine.createSpy('open').and.returnValue({ afterClosed: () => of({ address: 'addr1', amount: '3', hours: '9' }) });
    component.useLink();
    expect(component.destControls[0].value).toEqual({ address: 'addr1', coins: '3', hours: '9' });
    expect(component.autoHours).toBeFalse();
  });
});
