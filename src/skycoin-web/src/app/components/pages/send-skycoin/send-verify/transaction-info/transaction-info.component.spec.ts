import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { TransactionInfoComponent } from './transaction-info.component';
import { PriceService } from '../../../../../services/price.service';
import { CoinService } from '../../../../../services/coin.service';
import { MockTranslatePipe, MockPriceService, MockCoinService, MockDateTimePipe, MockMsgBarService } from '../../../../../utils/test-mocks';
import { MsgBarService } from '../../../../../services/msg-bar.service';
import { getTxNote, setTxNote } from '../../../../../utils/tx-notes';

describe('TransactionInfoComponent', () => {
  let component: TransactionInfoComponent;
  let fixture: ComponentFixture<TransactionInfoComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ TransactionInfoComponent, MockTranslatePipe, MockDateTimePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: PriceService, useClass: MockPriceService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TransactionInfoComponent);
    component = fixture.componentInstance;
    component.transaction = { inputs: [], outputs: [] };
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('edits the note of a transaction from the history', () => {
    component.isPreview = false;
    component.transaction = { txid: 'note-txid', inputs: [], outputs: [] } as any;
    component.editNote();
    expect(component.editingNote).toBe(true);
    component.saveNote('  rent  ');
    expect(component.editingNote).toBe(false);
    expect(component.transaction.note).toBe('rent');
    expect(getTxNote('note-txid')).toBe('rent');
    setTxNote('note-txid', '');
  });
});
