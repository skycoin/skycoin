import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPassphraseDialogComponent } from './hw-passphrase-dialog.component';

describe('HwPassphraseDialogComponent', () => {
  let component: HwPassphraseDialogComponent;
  let fixture: ComponentFixture<HwPassphraseDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPassphraseDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPassphraseDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
