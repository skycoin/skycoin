import { Component } from '@angular/core';
import { WalletService } from './services/wallet.service';
import { BlockchainService } from './services/blockchain.service';
import 'rxjs/add/operator/takeWhile';
import { ApiService } from './services/api.service';
import { Http, Headers } from '@angular/http';
import {Observable} from 'rxjs/Observable';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  current: number;
  highest: number;
  version: string;
  releaseVersion: string;
  updateAvailable: boolean;

  constructor(
    public walletService: WalletService,
    private apiService: ApiService,
    private blockchainService: BlockchainService,
    private http: Http,
  ) {}

  ngOnInit() {
    this.setVersion();
  }

  private setVersion() {
    // Set build version
    this.apiService.get('version')
      .subscribe(output =>  this.version = output.version);

    // Set latest release version from github
    this.http.get('https://api.github.com/repos/skycoin/skycoin/tags')
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || 'Unable to fetch latest release version from github.'))
      .subscribe(response =>  {
        // Iterate though the tags
        // Find the latest tag which is not a rc
        for ( const i in response ) {
          if ( response.hasOwnProperty(i) && response[i]['name'].indexOf('rc') === -1 ) {
            this.releaseVersion = response[i]['name'].substr(1);
            break;
          }
        }

        // Check if build version and release version differ
        this.updateAvailable = (this.version !== this.releaseVersion);
      });
  }

}
