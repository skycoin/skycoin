import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Version } from '../app.datatypes';
import BigNumber from 'bignumber.js';

@Injectable()
export class AppService {
  error: number;
  version: Version;

  get burnRate() {
    return this.burnRateInternal;
  }

  private burnRateInternal = new BigNumber(0.5);

  constructor(
    private apiService: ApiService,
  ) {}

  testBackend() {
    this.apiService.get('health').subscribe(response => {
        this.version = response.version;
        this.burnRateInternal = new BigNumber(response.user_verify_transaction.burn_factor);
        if (!response.csrf_enabled) {
          this.error = 3;
        }
      },
      () => this.error = 2,
    );
  }
}
