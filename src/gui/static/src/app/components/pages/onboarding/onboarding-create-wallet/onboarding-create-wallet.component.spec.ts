import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OnboardingCreateWalletComponent } from './onboarding-create-wallet.component';

describe('OnboardingCreateWalletComponent', () => {
  let component: OnboardingCreateWalletComponent;
  let fixture: ComponentFixture<OnboardingCreateWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OnboardingCreateWalletComponent ],
    })
    .compileComponents();
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
