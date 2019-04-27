import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPassphraseActivationDialogComponent } from './hw-passphrase-activation-dialog.component';

describe('HwPassphraseActivationDialogComponent', () => {
  let component: HwPassphraseActivationDialogComponent;
  let fixture: ComponentFixture<HwPassphraseActivationDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPassphraseActivationDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPassphraseActivationDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
