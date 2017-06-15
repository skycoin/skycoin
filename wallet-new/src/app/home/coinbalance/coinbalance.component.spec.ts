import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CoinbalanceComponent } from './coinbalance.component';

describe('CoinbalanceComponent', () => {
  let component: CoinbalanceComponent;
  let fixture: ComponentFixture<CoinbalanceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CoinbalanceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CoinbalanceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
