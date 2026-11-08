import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule, UntypedFormBuilder } from '@angular/forms';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of, throwError } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { JointTransactionComponent } from './joint-transaction.component';
import { JointTransaction, JointTransactionService, minimumFee } from '../../../services/wallet/joint-transaction.service';
import { WalletService } from '../../../services/wallet/wallet.service';
import { SpendingService } from '../../../services/wallet/spending.service';
import { BalanceService } from '../../../services/wallet/balance.service';
import { ClipboardService } from '../../../services/clipboard.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { CoinService } from '../../../services/coin.service';
import { MockTranslatePipe } from '../../../utils/test-mocks';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';

describe('JointTransactionComponent', () => {
  let component: JointTransactionComponent;
  let fixture: ComponentFixture<JointTransactionComponent>;
  let joint: jasmine.SpyObj<JointTransactionService>;
  let spending: jasmine.SpyObj<SpendingService>;
  let msgBar: jasmine.SpyObj<MsgBarService>;

  const wallet: any = { label: 'mine', addresses: [{ address: 'A' }], seed: 'seed' };
  const ownOutput = { hash: 'own', address: 'A', coins: new BigNumber('0.001'), calculated_hours: new BigNumber(1112) };
  const theirOutput = { hash: 'theirs', address: 'B', coins: new BigNumber('0.001'), calculated_hours: new BigNumber(5) };

  function decoded(signed: boolean[]): JointTransaction {
    return {
      encoded: 'hex-' + signed.join('-'),
      txid: 'txid',
      inputs: [
        { index: 0, hash: 'own', address: 'A', coins: new BigNumber('0.001'), hours: new BigNumber(1112), signed: signed[0], mine: true },
        { index: 1, hash: 'theirs', address: 'B', coins: new BigNumber('0.001'), hours: new BigNumber(5), signed: signed[1], mine: false },
      ],
      outputs: [],
      inputHours: new BigNumber(1117),
      fee: new BigNumber(117),
      minimumFee: new BigNumber(112),
      signedCount: signed.filter(s => s).length,
      fullySigned: signed.every(s => s),
    };
  }

  function create(queryTx?: string) {
    joint = jasmine.createSpyObj('JointTransactionService', ['outputs', 'create', 'decode', 'sign', 'minimumFee']);
    joint.minimumFee.and.callFake((h: BigNumber) => minimumFee(h, new BigNumber(10)));
    joint.outputs.and.returnValue(of([theirOutput]));
    spending = jasmine.createSpyObj('SpendingService', ['getWalletUnspentOutputs', 'injectTransaction']);
    spending.getWalletUnspentOutputs.and.returnValue(of([ownOutput]));
    msgBar = jasmine.createSpyObj('MsgBarService', ['showDone']);
    TestBed.configureTestingModule({
      declarations: [JointTransactionComponent, MockTranslatePipe],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        UntypedFormBuilder,
        { provide: JointTransactionService, useValue: joint },
        { provide: WalletService, useValue: {
          currentWallets: of([wallet, { label: 'hw', addresses: [], isHardware: true }]),
          isServerWallet: () => false,
          ensureKeys: () => of(undefined),
        } },
        { provide: SpendingService, useValue: spending },
        { provide: BalanceService, useValue: { startGettingBalances: () => {} } },
        { provide: ClipboardService, useValue: { copy: () => Promise.resolve('') } },
        { provide: MsgBarService, useValue: msgBar },
        { provide: CustomMatDialogService, useValue: {} },
        { provide: CoinService, useValue: { currentCoin: of({ coinSymbol: 'SKY', hoursName: 'Coin Hours' }) } },
        { provide: ActivatedRoute, useValue: { snapshot: { queryParamMap: convertToParamMap(queryTx ? { tx: queryTx } : {}) } } },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(JointTransactionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  function fillPayment(vendorHours: string) {
    component.createForm.get('wallet')!.setValue(wallet);
    component.createForm.get('outputs')!.setValue([ownOutput]);
    component.createForm.get('theirs')!.setValue('theirs');
    component.lookupTheirOutputs();
    component.destinations.at(0).setValue({ address: 'A', coins: '0.001', hours: '0' });
    component.destinations.at(1).setValue({ address: 'B', coins: '0.001', hours: vendorHours });
  }

  it('offers only software wallets', () => {
    create();
    expect(component.wallets.map(w => w.label)).toEqual(['mine']);
  });

  it('loads the outputs of the chosen wallet', () => {
    create();
    component.createForm.get('wallet')!.setValue(wallet);
    expect(spending.getWalletUnspentOutputs).toHaveBeenCalledWith(wallet);
    expect(component.ownOutputs).toEqual([ownOutput]);
  });

  it('works out the fee of a coin hour payment and creates it', () => {
    create();
    fillPayment('1000');
    expect(joint.outputs).toHaveBeenCalledWith(['theirs']);
    expect(component.inputHours.toString()).toBe('1117');
    expect(component.fee.toString()).toBe('117');
    expect(component.minimumFee.toString()).toBe('112');
    expect(component.createProblem).toBe('');

    joint.create.and.returnValue(of('unsigned-hex'));
    component.create();
    expect(joint.create).toHaveBeenCalledWith(['own', 'theirs'], [
      { address: 'A', coins: new BigNumber('0.001'), hours: new BigNumber(0) },
      { address: 'B', coins: new BigNumber('0.001'), hours: new BigNumber(1000) },
    ], null);
    expect(component.created).toBe('unsigned-hex');
  });

  it('refuses a fee under the minimum', () => {
    create();
    fillPayment('1010');
    expect(component.createProblem).toBe('joint-tx.fee-too-low');
  });

  it('needs a change address when coins are left over', () => {
    create();
    fillPayment('1000');
    component.destinations.at(1).patchValue({ coins: '0.0005' });
    expect(component.createProblem).toBe('joint-tx.need-change');
    component.createForm.get('changeAddress')!.setValue('A');
    expect(component.createProblem).toBe('');
  });

  it('flags a counterparty output the node does not know', () => {
    create();
    joint.outputs.and.returnValue(of([]));
    fillPayment('1000');
    expect(component.missingHashes).toEqual(['theirs']);
    expect(component.createProblem).toBe('joint-tx.unknown-outputs');
  });

  it('reviews, signs its inputs and broadcasts once fully signed', () => {
    create('pasted-hex');
    expect(component.mode).toBe(DoubleButtonActive.RightButton);
    component.signForm.get('wallet')!.setValue(wallet);
    joint.decode.and.returnValues(of(decoded([false, true])), of(decoded([true, true])));
    component.review();
    expect(joint.decode).toHaveBeenCalledWith('pasted-hex', wallet);
    expect(component.unsignedOwnInputs).toBe(1);

    joint.sign.and.returnValue(of('signed-hex'));
    component.sign();
    expect(joint.sign).toHaveBeenCalledWith(wallet, undefined, decoded([false, true]));
    expect(joint.decode).toHaveBeenCalledWith('signed-hex', wallet);
    expect(component.tx!.fullySigned).toBeTrue();
    expect(component.signForm.value.rawTx).toBe('hex-true-true');

    spending.injectTransaction.and.returnValue(of('txid'));
    component.broadcast();
    expect(spending.injectTransaction).toHaveBeenCalledWith('hex-true-true');
    expect(msgBar.showDone).toHaveBeenCalledWith('joint-tx.broadcast-done');
    expect(component.tx).toBeUndefined();
  });

  it('reviews again with the newly chosen wallet', () => {
    create('pasted-hex');
    const other: any = { label: 'other', addresses: [] };
    component.signForm.get('wallet')!.setValue(wallet);
    joint.decode.and.returnValue(of(decoded([false, false])));
    component.review();
    component.signForm.get('wallet')!.setValue(other);
    expect(joint.decode.calls.mostRecent().args).toEqual(['pasted-hex', other]);
  });

  it('explains when the wallet has nothing to sign', () => {
    create('pasted-hex');
    component.signForm.get('wallet')!.setValue(wallet);
    const notMine = decoded([false, false]);
    notMine.inputs.forEach(i => i.mine = false);
    joint.decode.and.returnValue(of(notMine));
    component.review();
    expect(component.signNote).toBe('joint-tx.owns-none');

    joint.decode.and.returnValue(of(decoded([true, false])));
    component.review();
    expect(component.signNote).toBe('joint-tx.own-signed');

    joint.decode.and.returnValue(of(decoded([false, false])));
    component.review();
    expect(component.signNote).toBe('');
  });

  it('shows why a review failed', () => {
    create('bad');
    joint.decode.and.returnValue(throwError(() => new Error('joint-tx.invalid-tx')));
    component.review();
    expect(component.error).toBe('joint-tx.invalid-tx');
    expect(component.tx).toBeUndefined();
  });

  it('opens a created transaction for co-signing', () => {
    create();
    fillPayment('1000');
    joint.create.and.returnValue(of('unsigned-hex'));
    joint.decode.and.returnValue(of(decoded([false, false])));
    component.create();
    component.cosignCreated();
    expect(component.mode).toBe(DoubleButtonActive.RightButton);
    expect(joint.decode).toHaveBeenCalledWith('unsigned-hex', wallet);
  });
});
