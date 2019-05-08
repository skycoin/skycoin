import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExchangeCreateComponent } from './exchange-create.component';

describe('ExchangeCreateComponent', () => {
  let component: ExchangeCreateComponent;
  let fixture: ComponentFixture<ExchangeCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ExchangeCreateComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExchangeCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
