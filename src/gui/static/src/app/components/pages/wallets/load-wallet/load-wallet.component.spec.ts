import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LoadWalletComponent } from './load-wallet.component';

describe('LoadWalletComponent', () => {
  let component: LoadWalletComponent;
  let fixture: ComponentFixture<LoadWalletComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoadWalletComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
