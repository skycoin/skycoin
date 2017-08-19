import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { ConnectionModel } from '../models/connection.model';
import { Observable } from 'rxjs/Observable';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import 'rxjs/add/operator/mergeMap';

@Injectable()
export class NetworkService {

  private automaticPeers: Subject<ConnectionModel[]> = new BehaviorSubject<ConnectionModel[]>([]);

  constructor(
    private apiService: ApiService,
  ) {
    this.loadData();
  }

  automatic(): Observable<ConnectionModel[]> {
    return this.automaticPeers.asObservable();
  }

  retrieveDefaultConnections(): Observable<ConnectionModel[]> {
    return this.apiService.post('network/defaultConnections')
      .map(output => output.map((address, index) => ({
        id: index + 1,
        address: address,
        listen_port: 6000,
      })));
  }

  private loadData(): void {
    this.retrieveConnections().subscribe(connections => this.automaticPeers.next(connections));
    IntervalObservable
      .create(5000)
      .flatMap(() => this.retrieveConnections())
      .subscribe(connections => this.automaticPeers.next(connections));
  }

  private retrieveConnections(): Observable<ConnectionModel[]> {
    return this.apiService.post('network/connections')
      .map(response => response.connections.sort((a, b) =>  a.id - b.id));
  }
}
