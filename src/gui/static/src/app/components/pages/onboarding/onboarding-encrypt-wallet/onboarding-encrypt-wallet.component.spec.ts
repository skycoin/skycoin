import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OnboardingEncryptWalletComponent } from './onboarding-encrypt-wallet.component';

describe('OnboardingEncryptWalletComponent', () => {
  let component: OnboardingEncryptWalletComponent;
  let fixture: ComponentFixture<OnboardingEncryptWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OnboardingEncryptWalletComponent ],
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
