import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';

import { TranslateService } from '@ngx-translate/core';

import { PreviewTransaction } from '../../../../app.datatypes';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { parseResponseMessage } from './../../../../utils/errors';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { OfflineTxComponent, OfflineTxMode } from '../offline-tx/offline-tx.component';
import { setTxNote } from '../../../../utils/tx-notes';

@Component({
    selector: 'app-send-verify',
    templateUrl: './send-verify.component.html',
    styleUrls: ['./send-verify.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class SendVerifyComponent implements OnDestroy {
  @ViewChild('sendButton') sendButton!: ButtonComponent;
  @ViewChild('backButton') backButton!: ButtonComponent;
  @Input() transaction!: PreviewTransaction;
  @Output() onBack = new EventEmitter<boolean>();

  constructor(
    private balanceService: BalanceService,
    private spendingService: SpendingService,
    private msgBarService: MsgBarService,
    private dialog: CustomMatDialogService,
    private changeDetectorRef: ChangeDetectorRef,
    private translate: TranslateService,
  ) {}

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  send() {
    this.msgBarService.hide();
    this.sendButton.resetState();
    this.sendButton.setLoading();
    this.backButton.setDisabled();

    this.spendingService.injectTransactions(this.transaction)
      .subscribe(
        txids => {
          this.onSuccess(txids);
          this.changeDetectorRef.markForCheck();
        },
        (error) => {
          this.onError(error);
          this.changeDetectorRef.markForCheck();
        }
      );
  }

  get rawTx(): string | undefined {
    return this.transaction && !this.parts ? this.transaction.encoded : undefined;
  }

  /** How many transactions the send takes, when it is more than one. */
  get parts(): number {
    return this.transaction && this.transaction.parts && this.transaction.parts.length > 1 ? this.transaction.parts.length : 0;
  }

  /** Shows the signed transaction to broadcast elsewhere, without sending it. */
  showRaw() {
    OfflineTxComponent.open(this.dialog, { mode: OfflineTxMode.Show, rawTx: this.rawTx, unsigned: false });
  }

  back() {
    this.onBack.emit(false);
  }

  private onSuccess(txids: string[]) {
    const noteSaved = !this.transaction.note || txids.every(txid => setTxNote(txid, this.transaction.note!));
    setTimeout(() => {
      if (noteSaved && txids.length > 1) {
        this.msgBarService.showDone(this.translate.instant('send.split-sent', { count: txids.length, txids: txids.join(', ') }));
      } else if (noteSaved) {
        this.msgBarService.showDone('send.sent');
      } else {
        this.msgBarService.showWarning('send.saving-note-error');
      }
      this.changeDetectorRef.markForCheck();
    });
    this.balanceService.startGettingBalances();
    this.onBack.emit(true);
  }

  private onError(error: any) {
    const errorMessage = error.message ? error.message : parseResponseMessage(error['_body']);
    this.msgBarService.showError(errorMessage);
    this.sendButton.resetState();
  }
}
