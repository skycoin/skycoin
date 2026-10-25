import { Injectable } from '@angular/core';
import { Observable, catchError, forkJoin, map, of, switchMap, timer } from 'rxjs';

import { ApiService } from './api.service';

/** Where the node learned of a peer. */
export enum ConnectionSources {
  Default = 'default',
  Exchange = 'exchange',
}

/** A peer the current coin's node is connected to. */
export interface Connection {
  address: string;
  listenPort: number;
  outgoing: boolean;
  height: number;
  lastSent: number;
  lastReceived: number;
  source: ConnectionSources;
}

/** The current coin's node's peers, for the network page. */
@Injectable()
export class NetworkService {
  private readonly updatePeriod = 5 * 1000;

  constructor(private apiService: ApiService) { }

  /**
   * The node's peers, asked for now and every few seconds while subscribed. A poll
   * that fails emits null and the polling goes on.
   */
  connections(): Observable<Connection[] | null> {
    return timer(0, this.updatePeriod).pipe(
      switchMap(() => forkJoin([
        this.apiService.get('network/defaultConnections'),
        this.apiService.get('network/connections'),
      ]).pipe(catchError(() => of(null)))),
      map((answer: [string[], any] | null) => answer === null ? null : ((answer[1] && answer[1].connections) || [])
        .map((c: any): Connection => ({
          address: c.address,
          listenPort: c.listen_port,
          outgoing: c.outgoing,
          height: c.height,
          lastSent: c.last_sent,
          lastReceived: c.last_received,
          source: (answer[0] || []).includes(c.address) ? ConnectionSources.Default : ConnectionSources.Exchange,
        }))
        .sort((a: Connection, b: Connection) => a.address.localeCompare(b.address))));
  }
}
