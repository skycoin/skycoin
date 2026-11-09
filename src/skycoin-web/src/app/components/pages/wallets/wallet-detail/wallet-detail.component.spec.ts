import { NEVER, of, throwError } from 'rxjs';
import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
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
});
