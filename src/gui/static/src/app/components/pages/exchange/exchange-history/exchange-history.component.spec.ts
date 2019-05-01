import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExchangeHistoryComponent } from './exchange-history.component';

describe('ExchangeHistoryComponent', () => {
  let component: ExchangeHistoryComponent;
  let fixture: ComponentFixture<ExchangeHistoryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ExchangeHistoryComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExchangeHistoryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
