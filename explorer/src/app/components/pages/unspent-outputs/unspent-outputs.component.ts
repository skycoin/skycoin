import { switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { BigNumber } from 'bignumber.js';
import { of, Subscription } from 'rxjs';

import { ApiService } from '../../../services/api/api.service';
import { GetCurrentBalanceResponse } from '../../../app.datatypes';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Different states of the list.
 */
enum ShowMoreStatus {
  /**
   * Only some elements are being shown.
   */
  showMore = 0,
  /**
   * Preparing to show all elements.
   */
  loading = 1,
  /**
   * All elements are being shown.
   */
  dontShowMore = 2,
}

/**
 * Page for showing the list of unspent outputs of a specific address.
 */
@Component({
  selector: 'app-unspent-outputs',
  templateUrl: './unspent-outputs.component.html',
  styleUrls: ['./unspent-outputs.component.scss']
})
export class UnspentOutputsComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerDataResponseKey = 'serv-dta-response';

  /**
   * Indicates how many outputs can be shown initially. If the address has more outputs, the
   * user will have to click a link for showing the rest. This helps to mantain good performance.
   */
  private readonly maxInitialElements = 100;

  /**
   * If set to true, makes the control stop responding to mouse clicks, to avoid problems with
   * slow operations.
   */
  disableClicks = false;
  /**
   * Allows to access the ShowMoreStatus enum in the HTML file.
   */
  showMoreStatus = ShowMoreStatus;
  /**
   * How many outputs the current address has.
   */
  totalOutputs = 0;
  /**
   * Current state of the output list.
   */
  showMoreOutputs = ShowMoreStatus.dontShowMore;
  /**
   * Current address.
   */
  address: string;
  /**
   * Unspent outputs of the current address.
   */
  outputs: GetCurrentBalanceResponse;
  /**
   * Outputs that will be shown in the UI.
   */
  outputsToShow: any[] = [];
  /**
   * Total number of coins in the unspent outputs.
   */
  coins: BigNumber = null;
  /**
   * Total number of hours in the unspent outputs.
   */
  hours: BigNumber = null;
  /**
   * Small text to be shown in variaous parts (not in the loading control) while loading the data.
   * It may also contain small error messages.
   */
  loadingMsg = 'general.loadingMsg';
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
    private route: ActivatedRoute,
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

    let savedData;

    // Get the URL params.
    this.pageSubscriptions.push(this.route.params.pipe(switchMap((params: Params) => {
      // Get the address and request the data.
      this.address = params['address'];

      // Get the data.
      // Use saved data or get from the server. If there is no saved data, savedData is null.
      savedData = checkSavedData ? this.getLocalValue(this.persistentServerDataResponseKey) : null;
      if (savedData) {
        oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
        return of(savedData.value);
      } else {
        return this.api.getCurrentBalance(params['address']);
      }
    })).subscribe(response => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerDataResponseKey, response);
      }

      this.outputs = response;
      this.totalOutputs = response.head_outputs.length;

      // Calculate the total number of coins and hours.
      this.coins = new BigNumber(0);
      this.hours = new BigNumber(0);
      response.head_outputs.map(o => {
        this.coins = this.coins.plus(o.coins);
        this.hours = this.hours.plus(o.calculated_hours);
      });

      if (this.outputs.head_outputs.length > this.maxInitialElements) {
        // Show only the max initial number of elements.
        this.outputsToShow = this.outputs.head_outputs.slice(0, this.maxInitialElements);
        // Indicate that there are additional outputs to show.
        this.showMoreOutputs = ShowMoreStatus.showMore;
      } else {
        this.outputsToShow = this.outputs.head_outputs;
      }

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadData(false);
      }
    }, error => {
      if (this.outputs !== undefined) {
        if (error.status >= 400 && error.status < 500) {
          // The address was not found.
          this.loadingMsg = 'general.noData';
          this.longErrorMsg = 'unspentOutputs.withoutOutputs';
        } else {
          // Error loading the data.
          this.loadingMsg = 'general.shortLoadingErrorMsg';
          this.longErrorMsg = 'general.longLoadingErrorMsg';
        }
      }
    }));

  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
  }

  /**
   * Makes the control show all the outputs.
   */
  showAll() {
    if (this.showMoreOutputs === ShowMoreStatus.showMore) {
      // If the process takes too long to be completed, the UI may end blocked and all mouse
      // clicks would be processed after finishing, which could make users think that the
      // application behaves erratically (specially if the user starts clicking
      // indiscriminately trying to make the app respond and that causes a navigation after
      // finishing the operation). To avoid this, disableClicks is set to true to make the
      // control ignore all mouse clicks temporarily.
      this.disableClicks = true;
      // Indicate that the elements are being loaded.
      this.showMoreOutputs = ShowMoreStatus.loading;

      // Load all the elements after 2 frames, to give the application time for updating the
      // UI, in case it gets blocked.
      setTimeout(() => {
        // Updates the list with the elements.
        this.outputsToShow = this.outputs.head_outputs;
        // Updates the UI after 2 frames.
        setTimeout(() => {
          // Idicate that all elements are being shown.
          this.showMoreOutputs = ShowMoreStatus.dontShowMore;
          // Accept mouse click.
          this.disableClicks = false;
        });
      }, 32);
    }
  }
}
