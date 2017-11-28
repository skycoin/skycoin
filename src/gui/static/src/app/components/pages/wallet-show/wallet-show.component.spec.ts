import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WalletShowComponent } from './wallet-show.component';

describe('WalletShowComponent', () => {
  let component: WalletShowComponent;
  let fixture: ComponentFixture<WalletShowComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletShowComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletShowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
