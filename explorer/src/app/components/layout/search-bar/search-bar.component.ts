import { Component, ElementRef, ViewChild, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription, of } from 'rxjs';
import { delay, mergeMap } from 'rxjs/operators';

import { SearchService, SearchError } from '../../../services/search/search.service';
import { ApiService } from '../../../services/api/api.service';

/**
 * Shows the search bar. Is the user search something, this component takes cares of finding it
 * and redirecting the user.
 */
@Component({
  selector: 'app-search-bar',
  templateUrl: './search-bar.component.html',
  styleUrls: ['./search-bar.component.scss']
})
export class SearchBarComponent implements OnDestroy {
  /**
   * Search field.
   */
  @ViewChild('input', { static: true }) input: ElementRef;

  /**
   * If true, it is not possible to start another search and the UI is shown busy.
   */
  searching = false;
  /**
   * Indicates if the node synchronization alert msg must be shown.
   */
  showSyncWarning = false;

  private nodeUrlSubscription: Subscription;
  private syncSubscription: Subscription;
  private operationSubscription: Subscription;

  constructor(
    public searchService: SearchService,
    private router: Router,
    private translate: TranslateService,
    private apiService: ApiService,
  ) {
    // Each time the node URL is changed, check if the node is in sync.
    this.nodeUrlSubscription = this.apiService.localNodeUrl.subscribe(() => {
      this.checkSyncState(0);
    });
  }

  ngOnDestroy() {
    if (this.nodeUrlSubscription) {
      this.nodeUrlSubscription.unsubscribe();
    }
    if (this.syncSubscription) {
      this.syncSubscription.unsubscribe();
    }
    if (this.operationSubscription && !this.operationSubscription.closed) {
      this.operationSubscription.unsubscribe();
    }
  }

  /**
   * Checks if the node is in sync. If it is not, the procedure will be automatically
   * repeated after a delay.
   *
   * @param delayMs Delay before performing the actual check.
   */
  private checkSyncState(delayMs: number) {
    if (this.syncSubscription) {
      this.syncSubscription.unsubscribe();
    }

    // Get the data.
    this.syncSubscription = of(0).pipe(delay(delayMs), mergeMap(() => this.apiService.getSyncState())).subscribe(response => {
      // The node is considered in sync if it is less than 3 blocks behind the top. The gap
      // in blocks is to avoid showing the alert if the operation is made just when a new block
      // has been detected but not added, which would happen very fast and should not merit
      // alerting the user.
      this.showSyncWarning = response.highest - response.current > 3;

      // If the node is not in sync, repeat the operation after some time.
      if (this.showSyncWarning) {
        this.checkSyncState(10000);
      }
    }, () => {
      if (this.showSyncWarning) {
        this.checkSyncState(10000);
      }
    });
  }

  /**
   * Starts a seach with the text entered in input.
   */
  search() {
    // Do not continue if a search is in progress.
    if (this.searching) {
      return;
    }

    const valueToSearch = this.input.nativeElement.value.trim();
    // Do not continue if the string is empty.
    if (valueToSearch.length < 1) {
      return;
    }

    this.searching = true;

    // Get where the user should be redirected to.
    const navCommands = this.searchService.processTerm(valueToSearch);
    if (navCommands.error) {
      if (navCommands.error === SearchError.invalidSearchTerm) {
        alert(this.translate.instant('search.unableToFind', { term: valueToSearch }));
      }
      this.searching = false;

      return;
    }

    // If there is no error, wait for the navigation commands needed for redirecting the user
    // to the page with the requested data.
    this.operationSubscription = navCommands.resultNavCommands.subscribe(
      result => {
        // Navigate to the page with the requested data.
        this.router.navigate(result);

        // Reset the control.
        this.searching = false;
        this.input.nativeElement.value = '';
      },
      () => this.searching = false
    );
  }
}
