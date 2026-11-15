import { waitForAsync, ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

import { SendVerifyComponent } from './send-verify.component';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { MockTranslatePipe, MockBalanceService, MockSpendingService, MockMsgBarService, MockCustomMatDialogService } from '../../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { of, throwError } from 'rxjs';
import { getTxNote, setTxNote } from '../../../../utils/tx-notes';
import { PreviewTransaction } from '../../../../app.datatypes';

describe('SendVerifyComponent', () => {
  let component: SendVerifyComponent;
  let fixture: ComponentFixture<SendVerifyComponent>;
  const button = { resetState: () => undefined, setLoading: () => undefined, setDisabled: () => undefined } as any;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ SendVerifyComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: SpendingService, useClass: MockSpendingService },
        { provide: MsgBarService, useClass: MockMsgBarService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
        { provide: TranslateService, useValue: { instant: (key: string, params: any) => key + ' ' + JSON.stringify(params) } },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendVerifyComponent);
    component = fixture.componentInstance;
    component.transaction = { inputs: [], outputs: [], from: '', to: '', encoded: '' };
    fixture.detectChanges();
    (TestBed.inject(BalanceService) as any).startGettingBalances = () => undefined;
    component.sendButton = button;
    component.backButton = button;
  });

  function splitSend(): PreviewTransaction {
    const parts = ['a', 'b'].map(e => Object.assign(new PreviewTransaction(), { inputs: [], outputs: [], encoded: e }));

    return Object.assign(new PreviewTransaction(), { inputs: [], outputs: [], from: '', to: [], encoded: '', parts });
  }

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('saves the note under the txid the node returns', () => {
    const spending = TestBed.inject(SpendingService) as any;
    spending.injectTransactions = () => of(['sent-txid']);
    component.transaction.note = 'groceries';
    component.send();
    expect(getTxNote('sent-txid')).toBe('groceries');
    setTxNote('sent-txid', '');
  });

  it('tells the user a send takes several transactions, and offers no single raw one', () => {
    component.transaction = splitSend();
    fixture.componentRef.setInput('transaction', component.transaction);
    fixture.detectChanges();

    expect(component.parts).toBe(2);
    expect(component.rawTx).toBeUndefined();
    expect(fixture.nativeElement.querySelector('.-split-notice')).toBeTruthy();
  });

  it('shows every txid once all parts are sent', fakeAsync(() => {
    const spending = TestBed.inject(SpendingService) as any;
    spending.injectTransactions = jasmine.createSpy().and.returnValue(of(['t1', 't2']));
    const msgBar = TestBed.inject(MsgBarService);
    spyOn(msgBar, 'showDone');
    component.transaction = splitSend();
    component.send();
    tick();

    expect(spending.injectTransactions).toHaveBeenCalledWith(component.transaction);
    expect(msgBar.showDone).toHaveBeenCalledWith('send.split-sent {"count":2,"txids":"t1, t2"}');
  }));

  it('shows which parts were sent when a later one fails', () => {
    const spending = TestBed.inject(SpendingService) as any;
    spending.injectTransactions = () => throwError(() => new Error('1 of 2 sent: t1'));
    const msgBar = TestBed.inject(MsgBarService);
    spyOn(msgBar, 'showError');
    component.transaction = splitSend();
    component.send();

    expect(msgBar.showError).toHaveBeenCalledWith('1 of 2 sent: t1');
  });
});
