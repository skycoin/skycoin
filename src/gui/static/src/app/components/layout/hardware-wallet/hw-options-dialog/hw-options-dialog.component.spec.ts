import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwOptionsDialogComponent } from './hw-options-dialog.component';

describe('HwOptionsDialogComponent', () => {
  let component: HwOptionsDialogComponent;
  let fixture: ComponentFixture<HwOptionsDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwOptionsDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwOptionsDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
