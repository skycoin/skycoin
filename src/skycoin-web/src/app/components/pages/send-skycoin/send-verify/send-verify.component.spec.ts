import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { SendVerifyComponent } from './send-verify.component';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { MockTranslatePipe, MockBalanceService, MockSpendingService, MockMsgBarService, MockCustomMatDialogService } from '../../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { of } from 'rxjs';
import { getTxNote, setTxNote } from '../../../../utils/tx-notes';

describe('SendVerifyComponent', () => {
  let component: SendVerifyComponent;
  let fixture: ComponentFixture<SendVerifyComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ SendVerifyComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: SpendingService, useClass: MockSpendingService },
        { provide: MsgBarService, useClass: MockMsgBarService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
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

  it('saves the note under the txid the node returns', () => {
    const spending = TestBed.inject(SpendingService) as any;
    spending.injectTransaction = () => of('sent-txid');
    (TestBed.inject(BalanceService) as any).startGettingBalances = () => undefined;
    const button = { resetState: () => undefined, setLoading: () => undefined, setDisabled: () => undefined } as any;
    component.sendButton = button;
    component.backButton = button;
    component.transaction.note = 'groceries';
    component.send();
    expect(getTxNote('sent-txid')).toBe('groceries');
    setTxNote('sent-txid', '');
  });
});
