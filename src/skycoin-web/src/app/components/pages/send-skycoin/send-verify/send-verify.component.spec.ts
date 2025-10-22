import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { SendVerifyComponent } from './send-verify.component';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { MockTranslatePipe, MockBalanceService, MockSpendingService, MockMsgBarService } from '../../../../utils/test-mocks';
import { MsgBarService } from '../../../../services/msg-bar.service';

describe('SendVerifyComponent', () => {
  let component: SendVerifyComponent;
  let fixture: ComponentFixture<SendVerifyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendVerifyComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: SpendingService, useClass: MockSpendingService },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendVerifyComponent);
    component = fixture.componentInstance;
    component.transaction = { inputs: [], outputs: [], from: '', to: '', encoded: '' };
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
