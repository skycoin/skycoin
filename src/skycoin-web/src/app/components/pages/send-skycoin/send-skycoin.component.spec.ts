import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { SendSkycoinComponent } from './send-skycoin.component';
import { MockTranslatePipe, MockCoinService, MockNavBarService } from '../../../utils/test-mocks';
import { CoinService } from '../../../services/coin.service';
import { NavBarService } from '../../../services/nav-bar.service';

describe('SendSkycoinComponent', () => {
  let component: SendSkycoinComponent;
  let fixture: ComponentFixture<SendSkycoinComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendSkycoinComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: NavBarService, useClass: MockNavBarService },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendSkycoinComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
