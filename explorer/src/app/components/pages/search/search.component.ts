import { first } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';

import { SearchService, SearchError } from '../../../services/search/search.service';
import { Subscription } from 'rxjs';

/**
 * Page for searching a term. It was not created to be opened during normal navigation, but just as
 * an intermediate page to be opened when using the browser search bar. It is opened by the browser
 * with the search term, checks what the final URL should be and redirects the user.
 */
@Component({
  templateUrl: './search.component.html',
})
export class SearchComponent implements OnInit, OnDestroy {

  searchTerm = '';
  /**
   * Error text that must be shown if there is a problem trying to find were the user must be
   * redirected to.
   */
  errorMsg: string;

  /**
   * Observable subscriptions that will be cleaned when closing the page.
   */
  private pageSubscriptions: Subscription[] = [];

  constructor(
    private searchService: SearchService,
    private route: ActivatedRoute,
    private router: Router,
  ) {}

  ngOnInit() {
    this.pageSubscriptions.push(this.route.params.pipe(first()).subscribe((params: Params) => {
      /**
       * Cancel the operation if there is no search term.
       */
      if (!params['term'] || (params['term'].trim() as string).length < 1) {
        this.errorMsg = 'search.unableToFind';
      }
      this.searchTerm = params['term'].trim();

      // Get where the user should be redirected to.
      const navCommands = this.searchService.processTerm(this.searchTerm);
      if (navCommands.error) {
        if (navCommands.error === SearchError.invalidSearchTerm) {
          this.errorMsg = 'search.unableToFind';
        }

        return;
      }

      // If there is no error, wait for the navigation commands needed for redirecting the user
      // to the page with the requested data.
      this.pageSubscriptions.push(navCommands.resultNavCommands.subscribe(
        result => {
          // Navigate and erase the current page from the browser history.
          this.router.navigate(result, { replaceUrl: true });
        },
        () => this.errorMsg = 'search.unableToFind'
      ));
    }));
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
  }
}
