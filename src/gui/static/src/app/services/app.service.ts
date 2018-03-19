import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs/Observable';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import { ConnectionModel } from '../models/connection.model';
import { Version } from '../app.datatypes';

@Injectable()
export class AppService {

  error: number;
  version: Version;

  constructor(
    private apiService: ApiService,
  ) {
    this.monitorConnections();
  }

  testBackend() {
    this.apiService.getVersion().first().subscribe(
      version => {
        this.version = version;
        this.apiService.getCsrf().subscribe(null, () => this.error = 3);
      }, () => this.error = 2
    );
  }

  private monitorConnections() {
    this.retrieveConnections().subscribe(connections => this.setConnectionError(connections));

    IntervalObservable
      .create(1500)
      .flatMap(() => this.retrieveConnections())
      .subscribe(connections => this.setConnectionError(connections));
  }

  private retrieveConnections(): Observable<ConnectionModel[]> {
    return this.apiService.get('network/connections');
  }

  private setConnectionError(response: any) {
    if (response.connections.length === 0) {
      this.error = 1;
    }
    if (response.connections.length > 0 && this.error === 1) {
      this.error = null;
    }
  }
}
