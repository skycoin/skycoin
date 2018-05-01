import { Component, OnInit } from '@angular/core';
import { NetworkService } from '../../../../services/network.service';

@Component({
  selector: 'app-network',
  templateUrl: './network.component.html',
  styleUrls: ['./network.component.scss']
})
export class NetworkComponent implements OnInit {

  defaultConnections;

  constructor(
    public networkService: NetworkService,
  ) { }

  ngOnInit() {
    this.networkService.retrieveDefaultConnections().first().subscribe(output => this.defaultConnections = output);
  }
}
