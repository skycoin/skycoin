import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwGenerateSeedDialogComponent } from './hw-generate-seed-dialog.component';

describe('HwGenerateSeedDialogComponent', () => {
  let component: HwGenerateSeedDialogComponent;
  let fixture: ComponentFixture<HwGenerateSeedDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwGenerateSeedDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwGenerateSeedDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
