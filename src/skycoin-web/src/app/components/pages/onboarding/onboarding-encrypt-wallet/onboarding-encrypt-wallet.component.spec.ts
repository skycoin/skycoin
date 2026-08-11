import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { UntypedFormBuilder } from '@angular/forms';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { OnboardingEncryptWalletComponent } from './onboarding-encrypt-wallet.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('OnboardingEncryptWalletComponent', () => {
  let component: OnboardingEncryptWalletComponent;
  let fixture: ComponentFixture<OnboardingEncryptWalletComponent>;

  beforeEach(waitForAsync(() => {
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
      providers: [ UntypedFormBuilder ],
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
