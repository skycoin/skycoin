import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwRestoreSeedDialogComponent } from './hw-restore-seed-dialog.component';

describe('HwRestoreSeedDialogComponent', () => {
  let component: HwRestoreSeedDialogComponent;
  let fixture: ComponentFixture<HwRestoreSeedDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwRestoreSeedDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwRestoreSeedDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
