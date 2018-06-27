import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WalletDetailComponent } from './wallet-detail.component';

describe('WalletDetailComponent', () => {
  let component: WalletDetailComponent;
  let fixture: ComponentFixture<WalletDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WalletDetailComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
