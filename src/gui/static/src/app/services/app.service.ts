import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Version } from '../app.datatypes';
import BigNumber from 'bignumber.js';
import { Http, Response } from '@angular/http';
import { shouldUpgradeVersion } from '../utils/semver';
import { AppConfig } from '../app.config';

@Injectable()
export class AppService {
  error: number;
  version: Version;
  fullCoinName = ' ';
  coinName = ' ';
  hoursName = ' ';
  hoursNameSingular = ' ';
  explorerUrl = ' ';

  get burnRate() {
    return this.burnRateInternal;
  }
  private burnRateInternal = new BigNumber(0.5);

  get updateAvailable(): boolean {
    return this.updateAvailableInternal;
  }
  private updateAvailableInternal = false;

  get lastestVersion(): string {
    return this.lastestVersionInternal;
  }
  private lastestVersionInternal = '';

  constructor(
    private apiService: ApiService,
    private http: Http,
  ) {}

  testBackend() {
    this.apiService.get('health').subscribe(response => {
        this.version = response.version;
        this.detectUpdateAvailable();
        this.burnRateInternal = new BigNumber(response.user_verify_transaction.burn_factor);

        this.fullCoinName = response.fiber.display_name;
        this.coinName = response.fiber.ticker;
        this.hoursName = response.fiber.coin_hours_display_name;
        this.hoursNameSingular = response.fiber.coin_hours_display_name_singular;
        this.explorerUrl = response.fiber.explorer_url;

        if (this.explorerUrl.endsWith('/')) {
          this.explorerUrl = this.explorerUrl.substr(0, this.explorerUrl.length - 1);
        }

        if (!response.csrf_enabled) {
          this.error = 3;
        }
      },
      () => this.error = 2,
    );
  }

  private detectUpdateAvailable() {
    if (AppConfig.urlForVersionChecking) {
      this.http.get(AppConfig.urlForVersionChecking)
        .retryWhen(errors => errors.delay(30000))
        .subscribe((response: Response) => {
          this.lastestVersionInternal = response.text().trim();
          if (this.lastestVersionInternal.startsWith('v')) {
            this.lastestVersionInternal = this.lastestVersionInternal.substr(1);
          }
          this.updateAvailableInternal = shouldUpgradeVersion(this.version.version, this.lastestVersionInternal);
        });
    }
  }
}
