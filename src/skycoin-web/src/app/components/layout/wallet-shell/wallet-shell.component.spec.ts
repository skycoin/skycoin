import { TestBed, waitForAsync, ComponentFixture } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA, Renderer2 } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Router } from '@angular/router';
import { of } from 'rxjs';

import { WalletShellComponent } from './wallet-shell.component';
import { MockLanguageService, MockTranslatePipe, MockTranslateService, MockCustomMatDialogService, MockMsgBarService, MockCoinService, MockHwWalletService } from '../../../utils/test-mocks';
import { LanguageService } from '../../../services/language.service';
import { CoinService } from '../../../services/coin.service';
import { CipherProvider, InitializationResults } from '../../../services/cipher.provider';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { Bip39WordListService } from '../../../services/bip39-word-list.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { HwWalletPinService } from '../../../services/hw-wallet-pin.service';

describe('WalletShellComponent', () => {
  let component: WalletShellComponent;
  let fixture: ComponentFixture<WalletShellComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletShellComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: HwWalletService, useClass: MockHwWalletService },
        // The component only stores a component reference on it.
        { provide: HwWalletPinService, useValue: {} },
        { provide: CoinService, useClass: MockCoinService },
        { provide: LanguageService, useClass: MockLanguageService },
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: Router, useValue: { events: of({}) } },
        { provide: CipherProvider, useValue: {
          initialize() { return of(InitializationResults.Ok); },
        } },
        { provide: Renderer2, useValue: { addClass: null, removeClass: null } },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        { provide: Bip39WordListService, useValue: {} },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletShellComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create the app', waitForAsync(() => {
    expect(component).toBeTruthy();
  }));
});
