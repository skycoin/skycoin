import { mergeMap, delay } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { Subject, BehaviorSubject, Observable, of, Subscription } from 'rxjs';

import { ApiService } from './api.service';

export enum ConnectionSources {
  /**
   * Default node to which the local node will always try to connect when started.
   */
  Default = 'default',
  /**
   * Informed by a remote node.
   */
  Exchange = 'exchange',
}

/**
 * Represents a connection the node local has with a remote node.
 */
export interface Connection {
  /**
   * Address of the remote node.
   */
  address: string;
  /**
   * Connection port.
   */
  listenPort: number;
  /**
   * If the connection is outgoing or not.
   */
  outgoing: boolean;
  /**
   * Highest block on the remote node.
   */
  height: number;
  /**
   * Last time in which data was sent to the remote node, in Unix time.
   */
  lastSent: number;
  /**
   * Last time in which data was received from the remote node, in Unix time.
   */
  lastReceived: number;
  /**
   * Source from were the remote node was discovered.
   */
  source: ConnectionSources;
}

/**
 * Allows to know if the local node is connected to any remote node to get the list of those nodes.
 */
@Injectable()
export class NetworkService {
  /**
   * Indicates if the local node is not currently connected to any remote node.
   */
  noConnections = false;

  /**
   * List of default addresses to which the local node will always try connect to when started.
   */
  private dataRefreshSubscription: Subscription;

  private trustedAddresses: string[];

  /**
   * Emits the lists of remote nodes to which the local node is currently connected.
   */
  private connectionsSubject: Subject<Connection[]> = new BehaviorSubject<Connection[]>([]);

  constructor(
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    // Start updating the data periodically.
    this.startDataRefreshSubscription(0);
  }

  /**
   * Gets the lists of remote nodes the local node is currently connected to.
   */
  connections(): Observable<Connection[]> {
    return this.connectionsSubject.asObservable();
  }

  /**
   * Makes the service start updating the data periodically. If this function was called
   * before, the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the balance.
   */
  private startDataRefreshSubscription(delayMs: number) {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.dataRefreshSubscription = of(0).pipe(delay(delayMs), mergeMap(() => {
        // Get the list of default remote nodes, but only if the list has not been
        // obtained before.
        if (!this.trustedAddresses) {
          return this.apiService.get('network/defaultConnections');
        } else {
          return of(this.trustedAddresses);
        }
      }), mergeMap(defaultConnectionsResponse => {
        this.trustedAddresses = defaultConnectionsResponse;

        // Get the list of current connections.
        return this.apiService.get('network/connections');
      })).subscribe(connectionsResponse => {
        if (connectionsResponse.connections === null || connectionsResponse.connections.length === 0) {
          this.noConnections = true;
          this.ngZone.run(() => this.connectionsSubject.next([]));

          return;
        }

        this.noConnections = false;

        // Process the obtained remote connections and convert them to a known object.
        const currentConnections = (connectionsResponse.connections as any[]).map<Connection>(connection => {
          return {
            address: connection.address,
            listenPort: connection.listen_port,
            outgoing: connection.outgoing,
            height: connection.height,
            lastSent: connection.last_sent,
            lastReceived: connection.last_received,
            source: this.trustedAddresses.find(trustedAddress => trustedAddress === connection.address) ? ConnectionSources.Default : ConnectionSources.Exchange,
          };
        }).sort((a, b) => a.address.localeCompare(b.address));

        this.ngZone.run(() => this.connectionsSubject.next(currentConnections));

        // Repeat the operation after an appropiate delay.
        this.startDataRefreshSubscription(5000);
      }, () => this.startDataRefreshSubscription(5000));
    });
  }
}
