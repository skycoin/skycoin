import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { WalletOptionsComponent } from './wallet-options.component';
import { MockTranslatePipe, MockCoinService, MockWalletService } from '../../../../../utils/test-mocks';
import { CoinService } from '../../../../../services/coin.service';
import { WalletService } from '../../../../../services/wallet/wallet.service';

describe('WalletOptionsComponent', () => {
  let component: WalletOptionsComponent;
  let fixture: ComponentFixture<WalletOptionsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletOptionsComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: WalletService, useClass: MockWalletService },
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: {} }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

describe('WalletOptionsComponent auto-sweep option', () => {
  function optionsFor(data: any): WalletOptionsComponent {
    TestBed.configureTestingModule({
      declarations: [ WalletOptionsComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: WalletService, useClass: MockWalletService },
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ]
    });
    const fixture = TestBed.createComponent(WalletOptionsComponent);
    fixture.detectChanges();

    return fixture.componentInstance;
  }

  it('is offered for a wallet the node keeps', () => {
    expect(optionsFor({ filename: 'a.wlt', walletType: 'bip44' }).showSweepOption).toBeTrue();
  });

  it('is offered for a hardware wallet, whose panel says why it is not swept', () => {
    expect(optionsFor({ filename: 'a.wlt', isHardware: true }).showSweepOption).toBeTrue();
  });

  it('is offered for a browser wallet, which the page sweeps', () => {
    expect(optionsFor({ label: 'browser' }).showSweepOption).toBeTrue();
  });

  it('is hidden for a temporary or xpub wallet', () => {
    expect(optionsFor({ filename: 't.wlt', temporary: true }).showSweepOption).toBeFalse();
    TestBed.resetTestingModule();
    expect(optionsFor({ filename: 'x.wlt', walletType: 'xpub' }).showSweepOption).toBeFalse();
  });
});
