import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs/Observable';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { Connection, Version } from '../app.datatypes';
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
    private ngZone: NgZone,
  ) {
    this.monitorConnections();
  }

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

  private monitorConnections() {
    this.retrieveConnections().subscribe(connections => this.setConnectionError(connections));

    this.ngZone.runOutsideAngular(() => {
      IntervalObservable
        .create(1500)
        .flatMap(() => this.retrieveConnections())
        .subscribe(connections => this.ngZone.run(() => {
          this.setConnectionError(connections);
        }));
    });
  }

  private retrieveConnections(): Observable<Connection[]> {
    return this.apiService.get('network/connections');
  }

  private setConnectionError(response: any) {
    if (response.connections === null || response.connections.length === 0) {
      this.error = 1;
    }
    if (response.connections !== null && response.connections.length > 0 && this.error === 1) {
      this.error = null;
    }
  }
}
