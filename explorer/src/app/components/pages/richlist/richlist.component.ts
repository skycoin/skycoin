import { first, map, switchMap } from 'rxjs';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Observable, of, Subscription } from 'rxjs';

import { ApiService } from '../../../services/api/api.service';
import { RichlistEntry } from '../../../app.datatypes';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Page with the list of the addresses with the most coins, a page at a time.
 */
@Component({
    templateUrl: './richlist.component.html',
    styleUrls: ['./richlist.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class RichlistComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerEntriesResponseKey = 'serv-rlp-response';

  /**
   * How many addresses are shown per page.
   */
  readonly pageSize = 50;

  /**
   * The addresses of the current page.
   */
  entries: RichlistEntry[] = [];
  longErrorMsg!: string;
  /**
   * The current page, starting at 0.
   */
  pageIndex = 0;
  /**
   * Whether there are more addresses after the current page.
   */
  hasNextPage = false;
  /**
   * How many addresses the whole list has, if the node says.
   */
  total?: number;

  private pageSubscriptions: Subscription[] = [];
  private dataSubscription?: Subscription;

  constructor(
    private api: ApiService,
    public explorer: ExplorerService,
    private route: ActivatedRoute,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    super();
  }

  ngOnInit() {
    this.pageSubscriptions.push(this.route.paramMap.subscribe(params => {
      const page = parseInt(params.get('page') || '1', 10);
      this.pageIndex = page > 0 ? page - 1 : 0;
      this.entries = [];
      this.longErrorMsg = '';
      this.changeDetectorRef.markForCheck();
      this.loadData(true);
    }));

    return super.ngOnInit();
  }

  /**
   * How many pages the list has, if the node says how long it is.
   */
  get pageCount(): number | undefined {
    return this.total !== undefined ? Math.ceil(this.total / this.pageSize) : undefined;
  }

  /**
   * Coin amounts with a fixed number of decimal places, so the decimal points line up.
   */
  get coinsFormat(): string {
    const decimals = this.explorer.maxDecimals || 0;

    return '1.' + decimals + '-' + decimals;
  }

  private loadData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;
    const key = this.persistentServerEntriesResponseKey + '-' + this.pageIndex;

    const start = this.pageIndex * this.pageSize;
    const savedData = checkSavedData ? this.getLocalValue(key) : null;
    // One address more than the page, so a node without offset paging still shows if page 1
    // has a next one.
    let nextOperation: Observable<RichlistPage> = this.api.getRichlist(start, this.pageSize + 1).pipe(
      switchMap(response => {
        if (response.total !== undefined || start === 0) {
          return of<RichlistPage>({
            entries: response.entries.slice(0, this.pageSize),
            total: response.total,
            more: response.entries.length > this.pageSize,
          });
        }

        // A node without offset paging returns the top addresses. Ask it for every address up
        // to the end of this page and one more, which tells if there is a next page.
        return this.api.getRichlist(0, start + this.pageSize + 1).pipe(map((all): RichlistPage => ({
          entries: all.entries.slice(start, start + this.pageSize),
          more: all.entries.length > start + this.pageSize,
        })));
      }),
    );
    if (savedData) {
      nextOperation = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    }

    this.dataSubscription?.unsubscribe();
    this.dataSubscription = nextOperation.pipe(first()).subscribe(page => {
      if (!savedData) {
        this.saveLocalValue(key, page);
      }
      this.entries = page.entries || [];
      this.total = page.total;
      this.hasNextPage = page.total !== undefined ? start + this.entries.length < page.total : !!page.more;
      if (!this.entries.length) {
        this.longErrorMsg = 'richlist.noEntries';
      }
      this.changeDetectorRef.markForCheck();

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadData(false);
      }
    },
      () => {
        if (!this.entries.length) {
          // Error loading the data.
          this.longErrorMsg = 'general.longLoadingErrorMsg';
          this.changeDetectorRef.markForCheck();
        }
    });
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
    this.dataSubscription?.unsubscribe();
  }
}

/**
 * One page of the richlist. total comes from nodes that page the list themselves, more from
 * the fallback for older nodes.
 */
interface RichlistPage {
  entries: RichlistEntry[];
  total?: number;
  more?: boolean;
}
