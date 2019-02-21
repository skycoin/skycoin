import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import 'rxjs/add/operator/mergeMap';
import { Connection } from '../app.datatypes';

@Injectable()
export class NetworkService {
  noConnections = false;

  private automaticPeers: Subject<Connection[]> = new BehaviorSubject<Connection[]>([]);

  constructor(
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    this.loadData();
  }

  automatic(): Observable<Connection[]> {
    return this.automaticPeers.asObservable();
  }

  retrieveDefaultConnections(): Observable<Connection[]> {
    return this.apiService.get('network/defaultConnections')
      .map(output => output.map((address, index) => ({
        id: index + 1,
        address: address,
        listen_port: 6000,
      })));
  }

  private loadData(): void {
    this.retrieveConnections().subscribe(connections => this.automaticPeers.next(connections));

    this.ngZone.runOutsideAngular(() => {
      IntervalObservable
        .create(5000)
        .flatMap(() => this.retrieveConnections())
        .subscribe(connections =>  this.ngZone.run(() => {
          this.automaticPeers.next(connections);
        }));
    });
  }

  private retrieveConnections(): Observable<Connection[]> {
    return this.apiService.get('network/connections')
      .map(response => {
        if (response.connections === null || response.connections.length === 0) {
          this.noConnections = true;

          return [];
        }

        this.noConnections = false;

        return response.connections.sort((a, b) =>  a.id - b.id);
      });
  }
}
