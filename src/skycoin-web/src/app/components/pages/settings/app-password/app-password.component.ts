import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ViewChild } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { ButtonComponent } from '../../../layout/button/button.component';

/**
 * Sets, changes or removes the password the wallets in this browser are stored
 * under. The password never leaves the browser.
 */
@Component({
  selector: 'app-app-password',
  templateUrl: './app-password.component.html',
  styleUrls: ['./app-password.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class AppPasswordComponent {
  @ViewChild('saveButton') saveButton!: ButtonComponent;
  form: UntypedFormGroup;
  error = '';

  constructor(
    formBuilder: UntypedFormBuilder,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.form = formBuilder.group({ current: [''], password: [''], confirm: [''] });
  }

  get hasPassword(): boolean {
    return this.walletService.hasAppPassword;
  }

  save() {
    const { current, password, confirm } = this.form.value;
    if (this.hasPassword && !this.walletService.checkAppPassword(current)) {
      return this.fail('app-password.wrong-current');
    }
    if (!password) {
      return this.fail('app-password.empty');
    }
    if (password !== confirm) {
      return this.fail('app-password.mismatch');
    }

    const changing = this.hasPassword;
    this.saveButton.setLoading();
    this.walletService.setAppPassword(password).then(() => {
      this.done(changing ? 'app-password.changed' : 'app-password.was-set');
    });
  }

  remove() {
    this.walletService.removeAppPassword(this.form.value.current).then(
      () => this.done('app-password.removed'),
      () => this.fail('app-password.wrong-current'),
    );
  }

  lock() {
    // The password and the decrypted wallets live only in this page's memory.
    window.location.reload();
  }

  private fail(error: string) {
    this.error = error;
    this.changeDetectorRef.markForCheck();
  }

  private done(message: string) {
    this.error = '';
    this.form.reset({ current: '', password: '', confirm: '' });
    if (this.saveButton) {
      this.saveButton.resetState();
    }
    this.msgBarService.showDone(message);
    this.changeDetectorRef.markForCheck();
  }
}
