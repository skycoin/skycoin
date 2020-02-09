import { delay, retryWhen } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { BigNumber } from 'bignumber.js';
import { HttpClient } from '@angular/common/http';

import { ApiService } from './api.service';
import { shouldUpgradeVersion } from '../utils/general-utils';
import { AppConfig } from '../app.config';
import { redirectToErrorPage } from '../utils/errors';

/**
 * Allows to access general information about the application and the node.
 */
@Injectable()
export class AppService {
  /**
   * Indicates if the csrf token protection is disabled on the local node.
   */
  get csrfDisabled() {
    return this.csrfDisabledInternal;
  }
  private csrfDisabledInternal = false;

  /**
   * Version number of the local node.
   */
  get nodeVersion() {
    return this.nodeVersionInternal;
  }
  private nodeVersionInternal: string;

  /**
   * Complete name of the coin.
   */
  get fullCoinName() {
    return this.fullCoinNameInternal;
  }
  private fullCoinNameInternal = ' ';

  /**
   * Short name of the coin, normally 3 letters.
   */
  get coinName() {
    return this.coinNameInternal;
  }
  private coinNameInternal = ' ';

  /**
   * Name of the coin hours, in plural.
   */
  get hoursName() {
    return this.hoursNameInternal;
  }
  private hoursNameInternal = ' ';

  /**
   * Name of the coin hours, in singular.
   */
  get hoursNameSingular() {
    return this.hoursNameSingularInternal;
  }
  private hoursNameSingularInternal = ' ';

  /**
   * URL for accessing the blockchain explorer.
   */
  get explorerUrl() {
    return this.explorerUrlInternal;
  }
  private explorerUrlInternal = ' ';

  /**
   * Indicates the maximum number of decimals for the coin the node currently accepts.
   */
  get currentMaxDecimals() {
    return this.currentMaxDecimalsInternal;
  }
  private currentMaxDecimalsInternal = 6;

  /**
   * Rate used for calculating the amount of hours that should be burned as transaction fee
   * when sending coins. The minimum amount to burn is "totalHours / burnRate".
   */
  get burnRate() {
    return this.burnRateInternal;
  }
  private burnRateInternal = new BigNumber(2);

  /**
   * Indicates if there is an update for this app available for download.
   */
  get updateAvailable(): boolean {
    return this.updateAvailableInternal;
  }
  private updateAvailableInternal = false;

  /**
   * Number of the lastest version of this app, obtained from a remote service.
   */
  get lastestVersion(): string {
    return this.lastestVersionInternal;
  }
  private lastestVersionInternal = '';

  constructor(
    private apiService: ApiService,
    private http: HttpClient,
  ) {}

  /**
   * Connects to the node to update all the data this service makes available. Should be
   * called when starting the app.
   */
  UpdateData() {
    this.apiService.get('health').subscribe(response => {
      this.nodeVersionInternal = response.version.version;
      this.burnRateInternal = new BigNumber(response.user_verify_transaction.burn_factor);
      this.currentMaxDecimalsInternal = response.user_verify_transaction.max_decimals;

      this.detectUpdateAvailable();

      this.fullCoinNameInternal = response.fiber.display_name;
      this.coinNameInternal = response.fiber.ticker;
      this.hoursNameInternal = response.fiber.coin_hours_display_name;
      this.hoursNameSingularInternal = response.fiber.coin_hours_display_name_singular;
      this.explorerUrlInternal = response.fiber.explorer_url;

      if (this.explorerUrlInternal.endsWith('/')) {
        this.explorerUrlInternal = this.explorerUrlInternal.substr(0, this.explorerUrl.length - 1);
      }

      if (!response.csrf_enabled) {
        this.csrfDisabledInternal = true;
      }
    }, () => redirectToErrorPage(2));
  }

  /**
   * Consults in a remote server the number of the lastest version and compares it to the
   * value of this.lastestVersionInternal. If there is an update available, it sets
   * this.updateAvailableInternal to true.
   */
  private detectUpdateAvailable() {
    if (AppConfig.urlForVersionChecking) {
      this.http.get(AppConfig.urlForVersionChecking, { responseType: 'text' })
        .pipe(retryWhen(errors => errors.pipe(delay(30000))))
        .subscribe((response: string) => {
          this.lastestVersionInternal = response.trim();
          if (this.lastestVersionInternal.startsWith('v')) {
            this.lastestVersionInternal = this.lastestVersionInternal.substr(1);
          }
          this.updateAvailableInternal = shouldUpgradeVersion(this.nodeVersionInternal, this.lastestVersionInternal);
        });
    }
  }
}
