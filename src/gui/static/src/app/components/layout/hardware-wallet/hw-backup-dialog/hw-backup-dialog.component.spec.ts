import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwBackupDialogComponent } from './hw-backup-dialog.component';

describe('HwBackupDialogComponent', () => {
  let component: HwBackupDialogComponent;
  let fixture: ComponentFixture<HwBackupDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwBackupDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwBackupDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
