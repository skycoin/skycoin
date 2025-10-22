import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatCheckboxModule, MatProgressSpinnerModule, MatIconModule, MatTooltipModule } from '@angular/material';
import { FormBuilder } from '@angular/forms';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { OnboardingEncryptWalletComponent } from './onboarding-encrypt-wallet.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('OnboardingEncryptWalletComponent', () => {
  let component: OnboardingEncryptWalletComponent;
  let fixture: ComponentFixture<OnboardingEncryptWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        OnboardingEncryptWalletComponent,
        ButtonComponent,
        MockTranslatePipe
      ],
      imports: [
        MatCheckboxModule,
        MatTooltipModule,
        MatIconModule,
        MatProgressSpinnerModule
      ],
      providers: [ FormBuilder ],
      schemas: [ NO_ERRORS_SCHEMA ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OnboardingEncryptWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
