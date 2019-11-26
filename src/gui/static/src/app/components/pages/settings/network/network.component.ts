import { Component, OnDestroy, OnInit } from '@angular/core';
import { NetworkService } from '../../../../services/network.service';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss'],
})
export class NetworkComponent implements OnInit, OnDestroy {
  peers: any;

  private subscription: ISubscription;

  constructor(
    public networkService: NetworkService,
  ) { }

  ngOnInit() {
    this.subscription = this.networkService.retrieveDefaultConnections().subscribe(trusted => {
      this.subscription = this.networkService.automatic().subscribe(peers => {
        this.peers = peers.map(peer => {
          peer.source = trusted.find(p => p.address === peer.address) ? 'default' : 'exchange';

          return peer;
        }).sort((a, b) => a.address.localeCompare(b.address));
      });
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}
