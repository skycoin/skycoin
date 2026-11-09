import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { UntypedFormBuilder } from '@angular/forms';
import { of } from 'rxjs';

import { CreateWalletFormComponent } from './create-wallet-form.component';
import { CoinService } from '../../../../../services/coin.service';
import { MockTranslatePipe, MockCoinService, MockNodeHealthService, MockCustomMatDialogService } from '../../../../../utils/test-mocks';
import { Bip39WordListService } from '../../../../../services/bip39-word-list.service';
import { NodeHealthService } from '../../../../../services/node-health.service';
import { CustomMatDialogService } from '../../../../../services/custom-mat-dialog.service';

describe('CreateWalletFormComponent', () => {
  let component: CreateWalletFormComponent;
  let fixture: ComponentFixture<CreateWalletFormComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ CreateWalletFormComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: NodeHealthService, useClass: MockNodeHealthService },
        UntypedFormBuilder,
        { provide: CoinService, useClass: MockCoinService },
        { provide: Bip39WordListService, useValue: { validateWord: true } },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService },
      ]
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

  it('flags a passphrase with a space at either end', () => {
    const cases: [string, boolean][] = [['', false], ['pass word', false], [' pass', true], ['pass ', true]];
    cases.forEach(([value, want]) => {
      component.form.get('seed_passphrase')!.setValue(value);
      expect(component.passphraseHasEdgeSpace).withContext(JSON.stringify(value)).toBe(want);
    });
  });

  describe('loading temporarily', () => {
    let dialog: any;
    const checkbox = { checked: true };

    beforeEach(() => {
      dialog = TestBed.inject(CustomMatDialogService);
      component.create = false;
    });

    function check(confirmed: boolean) {
      dialog.open = jasmine.createSpy('open').and.returnValue({ afterClosed: () => of(confirmed) });
      component.onTemporaryChange({ checked: true, source: checkbox } as any);
    }

    it('asks before loading a wallet temporarily', () => {
      check(true);

      expect(dialog.open.calls.argsFor(0)[1].data.text).toBe('wallet.new.temporary-warning');
      expect(component.loadTemporarily).toBe(true);
      expect(component.getData().temporary).toBe(true);
    });

    it('stays unchecked when the warning is declined', () => {
      check(false);

      expect(checkbox.checked).toBe(false);
      expect(component.loadTemporarily).toBe(false);
      expect(component.getData().temporary).toBe(false);
    });

    it('unchecks without asking', () => {
      check(true);
      dialog.open.calls.reset();
      component.onTemporaryChange({ checked: false, source: checkbox } as any);

      expect(dialog.open).not.toHaveBeenCalled();
      expect(component.getData().temporary).toBe(false);
    });

    it('never makes a new wallet temporary', () => {
      check(true);
      component.create = true;

      expect(component.getData().temporary).toBe(false);
    });
  });
});
