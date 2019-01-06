import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwSeedWordDialogComponent } from './hw-seed-word-dialog.component';

describe('HwSeedWordDialogComponent', () => {
  let component: HwSeedWordDialogComponent;
  let fixture: ComponentFixture<HwSeedWordDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwSeedWordDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwSeedWordDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
