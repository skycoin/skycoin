import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPassphraseHelpDialogComponent } from './hw-passphrase-help-dialog.component';

describe('HwPassphraseHelpDialogComponent', () => {
  let component: HwPassphraseHelpDialogComponent;
  let fixture: ComponentFixture<HwPassphraseHelpDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPassphraseHelpDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPassphraseHelpDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
