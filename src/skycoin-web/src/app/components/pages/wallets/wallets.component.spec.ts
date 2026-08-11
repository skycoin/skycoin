import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

import { WalletsComponent } from './wallets.component';
import { WalletService } from '../../../services/wallet/wallet.service';
import { CoinService } from '../../../services/coin.service';
import { MockTranslatePipe, MockCoinService, MockWalletService, MockTranslateService, MockCustomMatDialogService, MockHwWalletService, MockMsgBarService } from '../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { MsgBarService } from '../../../services/msg-bar.service';

describe('WalletsComponent', () => {
  let component: WalletsComponent;
  let fixture: ComponentFixture<WalletsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletsComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: MsgBarService, useClass: MockMsgBarService },
        { provide: HwWalletService, useClass: MockHwWalletService },
        { provide: WalletService, useClass: MockWalletService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});

