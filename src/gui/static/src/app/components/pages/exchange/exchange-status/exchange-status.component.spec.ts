import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExchangeStatusComponent } from './exchange-status.component';

describe('ExchangeStatusComponent', () => {
  let component: ExchangeStatusComponent;
  let fixture: ComponentFixture<ExchangeStatusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ExchangeStatusComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExchangeStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
