import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatMenuModule } from '@angular/material/menu';
import { RouterTestingModule } from '@angular/router/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

import { WalletDetailComponent } from './wallet-detail.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { MockTranslatePipe, MockWalletService, MockTranslateService, MockCustomMatDialogService, MockMsgBarService, MockHwWalletService, MockCoinService } from '../../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { CoinService } from '../../../../services/coin.service';

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
        { provide: HwWalletService, useClass: MockHwWalletService },
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
});
