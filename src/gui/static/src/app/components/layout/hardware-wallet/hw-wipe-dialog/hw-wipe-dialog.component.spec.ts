import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwWipeDialogComponent } from './hw-wipe-dialog.component';

describe('HwWipeDialogComponent', () => {
  let component: HwWipeDialogComponent;
  let fixture: ComponentFixture<HwWipeDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwWipeDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwWipeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
