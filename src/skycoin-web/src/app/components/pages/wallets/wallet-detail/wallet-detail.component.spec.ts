import { NEVER, of, throwError } from 'rxjs';
import { waitForAsync, ComponentFixture, TestBed, fakeAsync, tick, flush } from '@angular/core/testing';
import { BigNumber } from 'bignumber.js';
import { MatMenuModule } from '@angular/material/menu';
import { RouterTestingModule } from '@angular/router/testing';
import { ChangeDetectorRef, NO_ERRORS_SCHEMA } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

import { WalletDetailComponent } from './wallet-detail.component';
import { AddressExposureService } from '../../../../services/wallet/address-exposure.service';
import { WalletService } from '../../../../services/wallet/wallet.service';
import {
  MockTranslatePipe, MockWalletService, MockTranslateService, MockCustomMatDialogService, MockMsgBarService, MockHwWalletService, MockCoinService,
  MockBlockchainService } from '../../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { CoinService } from '../../../../services/coin.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { NumberOfAddressesComponent } from '../number-of-addresses/number-of-addresses.component';

describe('WalletDetailComponent', () => {
  let component: WalletDetailComponent;
  let fixture: ComponentFixture<WalletDetailComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletDetailComponent, MockTranslatePipe ],
      imports: [
        RouterTestingModule,
        MatMenuModule
      ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: AddressExposureService, useValue: { getExposed: () => of(new Set<string>()) } },
        { provide: HwWalletService, useClass: MockHwWalletService },
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: WalletService, useClass: MockWalletService },
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletDetailComponent);
    component = fixture.componentInstance;
    component.wallet = { label: '', addresses: [] };
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  describe('confirming a hardware wallet address', () => {
    let hw: any;
    let msgBar: MsgBarService;

    beforeEach(() => {
      hw = TestBed.inject(HwWalletService);
      msgBar = TestBed.inject(MsgBarService);
      (TestBed.inject(TranslateService) as any).instant = (key: string) => key;
      component.wallet = { label: 'hw', isHardware: true, addresses: [{ address: 'a0' }, { address: 'a1' }] };
      spyOn(hw, 'checkIfCorrectHwConnected').and.returnValue(of(true));
      hw.confirmAddress = jasmine.createSpy('confirmAddress').and.returnValue(of({}));
    });

    it('checks the device, then shows the address on it', () => {
      spyOn(msgBar, 'showDone');
      component.onConfirmAddress(1);

      expect(hw.checkIfCorrectHwConnected).toHaveBeenCalledWith('a0');
      expect(hw.confirmAddress).toHaveBeenCalledWith(1);
      expect(msgBar.showDone).toHaveBeenCalledWith('hardware-wallet.confirm-address.confirmation');
      expect(component.confirmingIndex).toBeNull();
    });

    it('reports a device error', () => {
      spyOn(msgBar, 'showError');
      hw.confirmAddress.and.returnValue(throwError(() => ({ translatableErrorMsg: 'hardware-wallet.errors.refused' })));
      component.onConfirmAddress(0);

      expect(msgBar.showError).toHaveBeenCalledWith('hardware-wallet.errors.refused');
      expect(component.confirmingIndex).toBeNull();
    });

    it('ignores a second request while one is on the device', () => {
      hw.confirmAddress.and.returnValue(NEVER);
      component.onConfirmAddress(0);
      component.onConfirmAddress(1);

      expect(hw.confirmAddress).toHaveBeenCalledTimes(1);
      expect(component.confirmingIndex).toBe(0);
    });
  });

  describe('scanning addresses', () => {
    function render() {
      fixture.componentRef.injector.get(ChangeDetectorRef).markForCheck();
      fixture.detectChanges();
    }

    it('shows a scan button on a wide window for a wallet that can scan', () => {
      const walletService = TestBed.inject(WalletService);
      spyOn(walletService, 'canFindUsedAddresses').and.returnValue(true);
      render();
      const button = fixture.nativeElement.querySelector('.btn-scan-addresses');

      expect(button).not.toBeNull();
      expect(button.classList).toContain('-not-on-small-and-below');
      expect(button.classList).not.toContain('-on-small-and-below-only');

      spyOn(component, 'onScanAddresses');
      button.click();
      expect(component.onScanAddresses).toHaveBeenCalled();
    });

    it('hides the scan button for a wallet that cannot scan', () => {
      render();

      expect(fixture.nativeElement.querySelector('.btn-scan-addresses')).toBeNull();
    });
  });

  describe('adding addresses', () => {
    let dialog: any;
    let walletService: any;

    beforeEach(() => {
      dialog = TestBed.inject(CustomMatDialogService);
      walletService = TestBed.inject(WalletService);
      walletService.addAddress = jasmine.createSpy('addAddress').and.returnValue(of(undefined));
      walletService.isBrowserBip44 = () => false;
      component.wallet = { label: 'w', seed: 'seed', nextSeed: 'next', addresses: [{ address: 'a0' }] };
    });

    function closingWith(...results: any[]) {
      dialog.open = jasmine.createSpy('open').and.returnValues(...results.map(r => ({ afterClosed: () => of(r) })));
    }

    it('asks how many and adds that many', fakeAsync(() => {
      closingWith(3);
      component.onAddNewAddress();
      tick();

      expect(dialog.open.calls.argsFor(0)[0]).toBe(NumberOfAddressesComponent);
      expect(walletService.addAddress).toHaveBeenCalledWith(component.wallet, true, undefined, undefined, 3);
      flush();
    }));

    it('adds nothing when the quantity dialog is cancelled', fakeAsync(() => {
      closingWith(undefined);
      component.onAddNewAddress();
      tick();

      expect(walletService.addAddress).not.toHaveBeenCalled();
      flush();
    }));

    it('warns before a gap of unused addresses a restore would not find', fakeAsync(() => {
      closingWith(20, true);
      component.onAddNewAddress();
      tick();

      expect(dialog.open.calls.count()).toBe(2);
      expect(dialog.open.calls.argsFor(1)[1].data.text).toBe('wallet.add-many-confirmation');
      expect(walletService.addAddress).toHaveBeenCalledWith(component.wallet, true, undefined, undefined, 20);
      flush();
    }));

    it('counts an address with coins as used', fakeAsync(() => {
      component.wallet.addresses = [{ address: 'a0' }, { address: 'a1', balance: new BigNumber(1) }];
      closingWith(4);
      component.onAddNewAddress();
      tick();

      expect(dialog.open.calls.count()).toBe(1);
      expect(walletService.addAddress).toHaveBeenCalledWith(component.wallet, true, undefined, undefined, 4);
      flush();
    }));

    it('does not add when the gap warning is declined', fakeAsync(() => {
      closingWith(25, false);
      component.onAddNewAddress();
      tick();

      expect(walletService.addAddress).not.toHaveBeenCalled();
      flush();
    }));

    it('adds one address to a hardware wallet without asking how many', () => {
      const hw: any = TestBed.inject(HwWalletService);
      hw.getAddresses = jasmine.createSpy('getAddresses').and.returnValue(of({ rawResponse: ['a1'] }));
      walletService.saveWallets = jasmine.createSpy('saveWallets');
      dialog.open = jasmine.createSpy('open');
      component.wallet = { label: 'hw', isHardware: true, addresses: [{ address: 'a0' }] };
      component.onAddNewAddress();

      expect(dialog.open).not.toHaveBeenCalled();
      expect(hw.getAddresses).toHaveBeenCalledWith(1, 1);
      expect(component.wallet.addresses.map(a => a.address)).toEqual(['a0', 'a1']);
    });
  });
});
