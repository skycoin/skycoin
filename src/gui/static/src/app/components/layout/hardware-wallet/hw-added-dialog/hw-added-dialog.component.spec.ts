import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwAddedDialogComponent } from './hw-added-dialog.component';

describe('HwAddedDialogComponent', () => {
  let component: HwAddedDialogComponent;
  let fixture: ComponentFixture<HwAddedDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwAddedDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwAddedDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
