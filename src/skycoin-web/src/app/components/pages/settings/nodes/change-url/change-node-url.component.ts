import { Component, Inject, OnInit, ViewChild, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import BigNumber from 'bignumber.js';
import { HttpClient } from '@angular/common/http';
import { TranslateService } from '@ngx-translate/core';

import { CoinService, TemporarilyAllowCoinResult } from '../../../../../services/coin.service';
import { ButtonComponent } from '../../../../layout/button/button.component';
import { isEqualOrSuperiorVersion } from '../../../../../utils/semver';
import { MsgBarService } from '../../../../../services/msg-bar.service';

@Component({
    selector: 'app-change-node-url',
    templateUrl: './change-node-url.component.html',
    styleUrls: ['./change-node-url.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class ChangeNodeURLComponent implements OnInit, OnDestroy {
  @ViewChild('action') actionButton: ButtonComponent;

  disableDismiss = false;
  showingUrlForm = true;
  form: UntypedFormGroup;

  nodeVersion: string;
  lastBlock: number;
  hoursBurnRate: string;
  coinName: string;

  // Known-good default servers offered in the dropdown. Selecting one just
  // populates the free-text field below (which stays editable); Verify keeps
  // working with whatever ends up in the field. Coin-aware: BTC gets a few
  // public ssl:// Electrum servers, other coins only the built-in default.
  nodeOptions: {label: string, value: string}[] = [];

  private newUrl: string;
  private verificationSubscription: Subscription;
  private initialURL: string;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: {coinId: number, url: string},
    public dialogRef: MatDialogRef<ChangeNodeURLComponent>,
    private formBuilder: UntypedFormBuilder,
    private coinService: CoinService,
    private http: HttpClient,
    private translate: TranslateService,
    private msgBarService: MsgBarService,
  ) {
    this.initialURL = data.url;
  }

  ngOnInit() {
    this.initForm();
    this.initNodeOptions();
  }

  ngOnDestroy() {
    this.removeVerificationSubscription();
    this.msgBarService.hide();
    this.coinService.removeTemporarilyAllowedCoin();
  }

  closePopup() {
    this.dialogRef.close();
  }

  startChange() {
    this.msgBarService.hide();

    this.newUrl = this.form.value.url.trim();

    if (this.newUrl !== '') {
      this.actionButton.setLoading();
      this.disableDismiss = true;

      if (this.newUrl.endsWith('/')) {
        this.newUrl = this.newUrl.substr(0, this.newUrl.length - 1);
      }

      const coinAllowed = this.coinService.temporarilyAllowCoin(this.data.coinId, this.newUrl);
      if (coinAllowed !== TemporarilyAllowCoinResult.OK) {
        setTimeout(() => {
          let text;
          if (coinAllowed === TemporarilyAllowCoinResult.AlreadyInUse) {
            text = this.translate.instant('nodes.change.url-error');
          } else {
            text = this.translate.instant('nodes.change.cancelled-error');
          }
          this.msgBarService.showError(text);
          this.actionButton.resetState();
          this.disableDismiss = false;
        }, 32);

        return;
      }

      this.removeVerificationSubscription();

      // Bitcoin's "node" is an ssl:// electrum server reached through the visor's
      // BTC gateway, which serves /v1/btc/health (tip height) — not the skycoin
      // /api/v1/health. Verify it there and skip the skycoin-specific version /
      // csrf / burn-rate checks.
      const coin = this.coinService.coins.find(c => c.id === this.data.coinId);
      if (coin && coin.isBitcoin()) {
        this.verificationSubscription = this.http.get(this.newUrl + '/v1/btc/health').subscribe((response: any) => {
          this.coinName = coin.coinName;
          this.nodeVersion = null;
          this.hoursBurnRate = null;
          this.lastBlock = (response && response.tip_height) || 0;

          this.actionButton.resetState();
          this.disableDismiss = false;
          this.showingUrlForm = false;
        }, () => this.cancelChange(false, false));

        return;
      }

      this.verificationSubscription = this.http.get(this.newUrl + '/api/v1/health').subscribe((response: any) => {
        this.nodeVersion = response.version.version;
        this.lastBlock = response.blockchain.head.seq;

        if (!isEqualOrSuperiorVersion(this.nodeVersion, '0.24.0')) {
          this.cancelChange(true, false);
          return;
        } else if (response.csrf_enabled) {
          this.cancelChange(false, true);
          return;
        }

        if (isEqualOrSuperiorVersion(this.nodeVersion, '0.25.0')) {
          this.hoursBurnRate = new BigNumber(100).dividedBy(response.user_verify_transaction.burn_factor).decimalPlaces(3, BigNumber.ROUND_FLOOR).toString() + '%';
          this.coinName = response.coin;
        } else {
          this.hoursBurnRate = '50%';
          this.coinName = null;
        }

        this.actionButton.resetState();
        this.disableDismiss = false;
        this.showingUrlForm = false;

      }, () => this.cancelChange(false, false));
    } else {
      this.completeChange();
    }
  }

  showUrlForm() {
    this.removeVerificationSubscription();
    this.showingUrlForm = true;
  }

  completeChange() {
    this.coinService.changeNodeUrl(this.data.coinId, this.newUrl);
    this.closePopup();

    if (this.initialURL !== this.newUrl) {
      setTimeout(() => this.msgBarService.showDone('nodes.change.url-changed'));
    }
  }

  private cancelChange(invalidNodeVersion: boolean, csrfActivated: boolean) {
    this.actionButton.resetState();
    this.disableDismiss = false;

    let errorMesasge: string;
    if (invalidNodeVersion) {
      errorMesasge = this.translate.instant('nodes.change.invalid-version-error');
    } else if (csrfActivated) {
      errorMesasge = this.translate.instant('nodes.change.csrf-error');
    } else {
      errorMesasge = this.translate.instant('nodes.change.connection-error');
    }

    this.msgBarService.showError(errorMesasge);
    this.actionButton.resetState();
  }

  private removeVerificationSubscription() {
    if (this.verificationSubscription) {
      this.verificationSubscription.unsubscribe();
    }
  }

  private initForm() {
    this.form = this.formBuilder.group({
      url: [this.data.url],
    });
  }

  // Writes the option chosen in the dropdown into the free-text field. The
  // field stays editable, so the user can tweak the value afterwards.
  onSelectPreset(value: string) {
    this.form.get('url').setValue(value);
  }

  // Builds the coin-aware list of default servers. The first entry always maps
  // to an empty value ("Default (built-in)") so the coin falls back to the
  // built-in node. For Bitcoin we add a few public ssl:// Electrum servers.
  private initNodeOptions() {
    const defaultOption = { label: this.translate.instant('nodes.change.preset-default'), value: '' };

    const coin = this.coinService.coins.find(c => c.id === this.data.coinId);
    if (coin && coin.isBitcoin()) {
      this.nodeOptions = [
        defaultOption,
        { label: 'ssl://electrum.blockstream.info:50002', value: 'ssl://electrum.blockstream.info:50002' },
        { label: 'ssl://fortress.qtornado.com:443', value: 'ssl://fortress.qtornado.com:443' },
        { label: 'ssl://electrum.emzy.de:50002', value: 'ssl://electrum.emzy.de:50002' },
      ];
    } else {
      this.nodeOptions = [defaultOption];
    }
  }
}
