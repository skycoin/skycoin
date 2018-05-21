import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OnboardingSafeguardComponent } from './onboarding-safeguard.component';

describe('OnboardingSafeguardComponent', () => {
  let component: OnboardingSafeguardComponent;
  let fixture: ComponentFixture<OnboardingSafeguardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OnboardingSafeguardComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OnboardingSafeguardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
