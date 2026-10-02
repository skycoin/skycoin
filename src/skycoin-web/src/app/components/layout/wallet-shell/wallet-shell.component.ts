import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit, Renderer2, ViewChild } from '@angular/core';
import { LanguageService } from '../../../services/language.service';
import { Router, NavigationEnd, Event } from '@angular/router';
import { Subscription, filter } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { config } from '../../../app.config';
import { CipherProvider, InitializationResults } from '../../../services/cipher.provider';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { Bip39WordListService } from '../../../services/bip39-word-list.service';
import { MsgBarComponent } from '../msg-bar/msg-bar.component';
import { MsgBarService } from '../../../services/msg-bar.service';
import { HwWalletPinService } from '../../../services/hw-wallet-pin.service';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { HwPinDialogComponent } from '../hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { HwConfirmTxDialogComponent } from '../hardware-wallet/hw-confirm-tx-dialog/hw-confirm-tx-dialog.component';
import { WALLET_SCOPE_CLASS } from '../../../wallet-scope';

/**
 * The wallet's frame: cipher start-up, the node status bar, the outlet for the
 * wallet's pages and the message bar. The standalone app shows only this; an
 * app hosting the wallet mounts it on a route of its own (see wallet.module.ts).
 */
@Component({
    selector: 'app-wallet-shell',
    templateUrl: './wallet-shell.component.html',
    styleUrls: ['./wallet-shell.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false,
    host: { class: WALLET_SCOPE_CLASS },
})
export class WalletShellComponent implements OnInit, OnDestroy {
  @ViewChild('msgBar') msgBar!: MsgBarComponent;

  current!: number;
  highest!: number;
  otcEnabled!: boolean;
  version!: string;
  browserCompatibleWithWasm = true;
  wasmFileLoaded = true;

  // What the shell subscribes to, released when it goes: an app hosting the
  // wallet on one of its routes creates a new shell on every visit.
  private subscriptions = new Subscription();

  constructor(
    private languageService: LanguageService,
    cipherProvider: CipherProvider,
    router: Router,
    dialog: CustomMatDialogService,
    private renderer: Renderer2,
    private bip38WordList: Bip39WordListService,
    private msgBarService: MsgBarService,
    hwWalletPinService: HwWalletPinService,
    hwWalletService: HwWalletService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    // Set component references to avoid circular dependencies
    hwWalletPinService.requestPinComponent = HwPinDialogComponent;
    hwWalletService.signTransactionConfirmationComponent = HwConfirmTxDialogComponent;
    this.subscriptions.add(router.events.pipe(
      filter((event: Event): event is NavigationEnd => event instanceof NavigationEnd)
    ).subscribe(() => {
      window.scrollTo(0, 0);
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptions.add(cipherProvider.initialize().subscribe(response => {
      this.checkCipherProviderResponse(response);
      this.changeDetectorRef.markForCheck();
    }, response => {
      this.checkCipherProviderResponse(response);
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptions.add(dialog.showingDialog.subscribe(value => {
      if (!value) {
        renderer.addClass(document.body, 'fix-error-position');
      } else {
        renderer.removeClass(document.body, 'fix-error-position');
      }
      this.changeDetectorRef.markForCheck();
    }));
  }

  ngOnInit() {
    this.otcEnabled = config.otcEnabled;
    this.languageService.loadLanguageSettings();

    this.msgBarService.msgBarComponent = this.msgBar;
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
    this.renderer.removeClass(document.body, 'fix-error-position');
  }

  loading() {
    return !this.current || !this.highest || this.current !== this.highest;
  }

  private checkCipherProviderResponse(response: any) {
    if ((window as any)['removeSplash']) {
      setTimeout(() => {
        (window as any)['removeSplash']();
        this.changeDetectorRef.markForCheck();
      });
    }
    if (response !== InitializationResults.Ok) {
      if (response === InitializationResults.ErrorLoadingWasmFile) {
        this.wasmFileLoaded = false;
      } else {
        this.browserCompatibleWithWasm = false;
      }
    }
  }
}
