import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubscriptionLike } from 'rxjs';

import { NetworkService, Connection } from '../../../../services/network.service';

/**
 * Allows to see the list of connections the node currently has with other nodes.
 */
@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss'],
})
export class NetworkComponent implements OnInit, OnDestroy {
  peers: Connection[];

  private subscription: SubscriptionLike;

  constructor(
    public networkService: NetworkService,
  ) { }

  ngOnInit() {
    // Periodically get the list of connected nodes.
    this.subscription = this.networkService.connections().subscribe(peers => this.peers = peers);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}
