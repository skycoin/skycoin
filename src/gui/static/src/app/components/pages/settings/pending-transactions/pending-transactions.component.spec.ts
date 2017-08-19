import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PendingTransactionsComponent } from './pending-transactions.component';

describe('PendingTransactionsComponent', () => {
  let component: PendingTransactionsComponent;
  let fixture: ComponentFixture<PendingTransactionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PendingTransactionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PendingTransactionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
