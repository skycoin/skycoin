import { first } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subscription } from 'rxjs';

import { ApiService } from '../../../services/api/api.service';
import { RichlistEntry } from '../../../app.datatypes';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Page for showing the list of unlocked addresses with most coins.
 */
@Component({
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.scss']
})
export class RichlistComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerEntriesResponseKey = 'serv-etr-response';

  /**
   * List to be shown in the UI.
   */
  entries: RichlistEntry[] = [];
  /**
   * Error message to be shown in the loading control if there is a problem.
   */
  longErrorMsg: string;

  /**
   * Observable subscriptions that will be cleaned when closing the page.
   */
  private pageSubscriptions: Subscription[] = [];

  constructor(
    private api: ApiService,
    public explorer: ExplorerService
  ) {
    super();
  }

  ngOnInit() {
    this.loadData(true);

    return super.ngOnInit();
  }

  private loadData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;

    // Get the data.
    // Use saved data or get from the server. If there is no saved data, savedData is null.
    const savedData = checkSavedData ? this.getLocalValue(this.persistentServerEntriesResponseKey) : null;
    let nextOperation: Observable<any> = this.api.getRichlist();
    if (savedData) {
      nextOperation = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    }

    this.pageSubscriptions.push(nextOperation.pipe(first()).subscribe(entries => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerEntriesResponseKey, entries);
      }
      this.entries = entries;

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadData(false);
      }
    },
      () => {
        if (!this.entries.length) {
          // Error loading the data.
          this.longErrorMsg = 'general.longLoadingErrorMsg';
        }
    }));
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
  }
}
