import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwSeedDialogComponent } from './hw-seed-dialog.component';

describe('HwSeedDialogComponent', () => {
  let component: HwSeedDialogComponent;
  let fixture: ComponentFixture<HwSeedDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwSeedDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwSeedDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
