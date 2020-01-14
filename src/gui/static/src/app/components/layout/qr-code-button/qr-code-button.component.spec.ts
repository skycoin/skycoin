import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { QrCodeButtonComponent } from './qr-code-button.component';

describe('QrCodeButtonComponent', () => {
  let component: QrCodeButtonComponent;
  let fixture: ComponentFixture<QrCodeButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ QrCodeButtonComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(QrCodeButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
