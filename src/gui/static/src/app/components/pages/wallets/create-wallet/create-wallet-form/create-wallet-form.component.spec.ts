import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { CreateWalletFormComponent } from './create-wallet-form.component';

describe('CreateWalletFormComponent', () => {
  let component: CreateWalletFormComponent;
  let fixture: ComponentFixture<CreateWalletFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateWalletFormComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateWalletFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
