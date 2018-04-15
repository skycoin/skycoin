import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OnboardingDisclaimerComponent } from './onboarding-disclaimer.component';

describe('OnboardingDisclaimerComponent', () => {
  let component: OnboardingDisclaimerComponent;
  let fixture: ComponentFixture<OnboardingDisclaimerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OnboardingDisclaimerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OnboardingDisclaimerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
