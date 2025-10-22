import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { FormBuilder } from '@angular/forms';

import { CreateWalletFormComponent } from './create-wallet-form.component';
import { CoinService } from '../../../../../services/coin.service';
import { MockTranslatePipe, MockCoinService } from '../../../../../utils/test-mocks';
import { Bip39WordListService } from '../../../../../services/bip39-word-list.service';

describe('CreateWalletFormComponent', () => {
  let component: CreateWalletFormComponent;
  let fixture: ComponentFixture<CreateWalletFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateWalletFormComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        FormBuilder,
        { provide: CoinService, useClass: MockCoinService },
        { provide: Bip39WordListService, useValue: { validateWord: true } }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateWalletFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
