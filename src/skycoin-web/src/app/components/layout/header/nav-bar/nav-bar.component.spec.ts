import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltip } from '@angular/material/tooltip';

import { NavBarComponent } from './nav-bar.component';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonComponent } from '../../double-button/double-button.component';
import { ButtonComponent } from '../../button/button.component';
import { MockTranslatePipe, MockNavBarService } from '../../../../utils/test-mocks';
import { WalletPathPipe } from '../../../../wallet-path';

describe('NavBarComponent', () => {
  let component: NavBarComponent;
  let fixture: ComponentFixture<NavBarComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ NavBarComponent, MockTranslatePipe, DoubleButtonComponent, ButtonComponent, WalletPathPipe ],
      imports: [ MatIconModule, MatProgressSpinnerModule, MatTooltip, RouterTestingModule ],
      providers: [ { provide: NavBarService, useClass: MockNavBarService } ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NavBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
