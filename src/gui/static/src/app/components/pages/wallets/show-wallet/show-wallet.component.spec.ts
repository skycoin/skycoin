import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowWalletComponent } from './show-wallet.component';

describe('ShowWalletComponent', () => {
  let component: ShowWalletComponent;
  let fixture: ComponentFixture<ShowWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShowWalletComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShowWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
