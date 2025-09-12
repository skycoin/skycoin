import { Subscription } from 'rxjs';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { ApiService } from 'app/services/api/api.service';
import { ExplorerService } from 'app/services/explorer/explorer.service';

/**
 * Page for setting the URL of the local node the app must use as backend. The new URL is saved
 * in local storage, so it is used again the next time the app is openned using the same
 * domain name. If the URL is the string "null", any previously saved URL is deleted. Just after
 * saving the URL, the browser is redirected to the main page.
 */
@Component({
  selector: 'app-node-url',
  template: '',
})
export class NodeUrlComponent implements OnInit, OnDestroy {
  // Subscriptions that will be cleaned when closing the page.
  private navParamsSubscription: Subscription;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private explorerService: ExplorerService,
    private apiService: ApiService,
  ) { }

  ngOnInit() {
    // Check the URL params to detect the URL of the local node.
    this.navParamsSubscription = this.route.params.subscribe(params => {
      let nodeUrl: string = params['url'];

      // If the url is "null", convert the value to a real null.
      if (nodeUrl.toUpperCase() === 'null'.toUpperCase()) {
        nodeUrl = null;
      }

      // Save and use the URL.
      this.apiService.setNodeUrl(nodeUrl);
      // Refresh the basic info about the backend.
      this.explorerService.initialize();
      // Redirect to the main page.
      this.router.navigate([''], { replaceUrl: true });
    });
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.navParamsSubscription.unsubscribe();
  }
}
