import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { SendSkycoinComponent } from './send-skycoin.component';
import { MockTranslatePipe, MockCoinService, MockNavBarService, MockCustomMatDialogService, MockMsgBarService } from '../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { CoinService } from '../../../services/coin.service';
import { NavBarService } from '../../../services/nav-bar.service';

describe('SendSkycoinComponent', () => {
  let component: SendSkycoinComponent;
  let fixture: ComponentFixture<SendSkycoinComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ SendSkycoinComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: NavBarService, useClass: MockNavBarService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        { provide: MsgBarService, useClass: MockMsgBarService },
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
