import { Component, OnDestroy, OnInit } from '@angular/core';
import { NetworkService, Connection } from '../../../../services/network.service';
import { SubscriptionLike } from 'rxjs';

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
    this.subscription = this.networkService.connections().subscribe(peers => this.peers = peers);
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}
