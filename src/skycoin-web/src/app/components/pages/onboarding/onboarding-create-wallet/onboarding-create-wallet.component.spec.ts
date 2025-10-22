import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { RouterTestingModule } from '@angular/router/testing';
import { FormBuilder } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { TranslateService } from '@ngx-translate/core';

import { OnboardingCreateWalletComponent } from './onboarding-create-wallet.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { CoinService } from '../../../../services/coin.service';
import { MockTranslatePipe, MockWalletService, MockCoinService, MockLanguageService, MockBlockchainService, MockCustomMatDialogService, MockMsgBarService } from '../../../../utils/test-mocks';
import { LanguageService } from '../../../../services/language.service';
import { CreateWalletFormComponent } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { BlockchainService } from '../../../../services/blockchain.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { Bip39WordListService } from '../../../../services/bip39-word-list.service';
import { MsgBarService } from '../../../../services/msg-bar.service';

describe('OnboardingCreateWalletComponent', () => {
  let component: OnboardingCreateWalletComponent;
  let fixture: ComponentFixture<OnboardingCreateWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        OnboardingCreateWalletComponent,
        MockTranslatePipe,
        CreateWalletFormComponent
      ],
      imports: [
        RouterTestingModule,
        BrowserAnimationsModule,
      ],
      providers: [
        FormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: LanguageService, useClass: MockLanguageService },
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        {
          provide: TranslateService,
          useValue: jasmine.createSpyObj('TranslateService', ['instant'])
        },
        { provide: Bip39WordListService, useValue: { validateWord: true } },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ],
      schemas: [ NO_ERRORS_SCHEMA ]
    });
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OnboardingCreateWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
