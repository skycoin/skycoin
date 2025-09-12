import { first, switchMap } from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Observable, of, Subscription } from 'rxjs';

import { ApiService } from '../../../services/api/api.service';
import { Block } from '../../../app.datatypes';
import { ExplorerService } from '../../../services/explorer/explorer.service';
import { PageBaseComponent } from '../page-base';
import { dataValidityTime } from 'app/app.config';

/**
 * Main page. Shows some details about the blockchain and a list of blocks.
 */
@Component({
  templateUrl: './blocks.component.html',
  styleUrls: ['./blocks.component.scss']
})
export class BlocksComponent extends PageBaseComponent implements OnInit, OnDestroy {
  // Keys for persisting the server data, to be able to restore the state after navigation.
  private readonly persistentServerMetadataResponseKey = 'serv-mtd-response';
  private readonly persistentServerBlocksResponseKey = 'serv-bks-response';
  private readonly persistentServerSupplyResponseKey = 'serv-spl-response';

  /**
   * List of blocks to show in the UI.
   */
  blocks: Block[] = [];
  /**
   * Current coin supply.
   */
  currentSupply: number;
  /**
   * Max coin supply.
   */
  totalSupply: number;
  /**
   * Current hour coin supply.
   */
  currentCoinhourSupply: number;
  /**
   * Max coin hour supply.
   */
  totalCoinhourSupply: number;
  /**
   * How many blocks the blockchain currently has.
   */
  blockCount = 0;
  /**
   * The current page.
   */
  pageIndex = 0;
  /**
   * How many blocks are shown per page.
   */
  readonly pageSize = 10;
  /**
   * Small text to be shown in the coin supply fields while loading the data.
   * It may also contain small error messages.
   */
  loadingCoinSupplyMsg = 'general.loadingMsg';
  /**
   * Small text to be shown in the block count field while loading the data.
   * It may also contain small error messages.
   */
  loadingMetadataMsg = 'general.loadingMsg';
  /**
   * Error message to be shown in the loading control if there is a problem.
   */
  longErrorMsg: string;
  /**
   * If the app is using a local node as backend.
   */
  usingLocalNode = false;

  /**
   * If the mouse is over the date of a block, this var contains the ID of that block,
   * for making the UI to show the timezone. The var is modified in the HTML file and
   * is used for small window sizes.
   */
  mouseOver = -1;

  /**
   * Observable subscriptions that will be cleaned when closing the page.
   */
  private nodeUrlSubscription: Subscription;
  private pageSubscriptions: Subscription[] = [];

  /**
   * The number of pages in which all the block are divided.
   */
  get pageCount() {
    return Math.ceil(this.blockCount / this.pageSize);
  }

  constructor(
    private api: ApiService,
    private explorer: ExplorerService,
    private route: ActivatedRoute,
  ) {
    super();

    // Each time the node URL is changed, check if a local node is being used.
    this.nodeUrlSubscription = this.api.localNodeUrl.subscribe(response => {
      this.usingLocalNode = !!response;
    });
  }

  ngOnInit() {
    this.loadBasicData(true);
    this.loadSupplyData(true);

    return super.ngOnInit();
  }

  private loadBasicData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;

    // Get how many blocks the blockchain has.
    // Use saved data or get from the server. If there is no saved data, savedData is null.
    let savedData = checkSavedData ? this.getLocalValue(this.persistentServerMetadataResponseKey) : null;
    let nextOperation: Observable<any> = this.api.getBlockchainMetadata();
    if (savedData) {
      nextOperation = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    }

    this.pageSubscriptions.push(nextOperation.pipe(first(), switchMap(blockchain => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerMetadataResponseKey, blockchain);
      }
      this.blockCount = blockchain.blocks;

      // Get the URL params.
      return this.route.paramMap.pipe(first());
    }), switchMap(params => {
      // Calculate the values of the current page.
      this.pageIndex = parseInt(params.get('page'), 10) - 1;
      const end = this.blockCount - (this.pageIndex * this.pageSize);
      const begin = end - this.pageSize + 1;

      // Request the blocks that will be displayed in the UI.
      // Use saved data or get from the server. If there is no saved data, savedData is null.
      savedData = checkSavedData ? this.getLocalValue(this.persistentServerBlocksResponseKey) : null;
      nextOperation = this.explorer.getBlocks(begin > 0 ? begin : 0, end > 0 ? end : 0).pipe(first());
      if (savedData) {
        nextOperation = of(savedData.value);
        oldSavedDataUsed = oldSavedDataUsed || (savedData.date < (new Date()).getTime() - dataValidityTime);
      }

      return nextOperation;
    })).subscribe(blocks => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerBlocksResponseKey, blocks);
      }
      this.blocks = blocks;

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadBasicData(false);
      }
    }, () => {
      if (!this.blocks.length) {
        // Error loading the data.
        this.loadingMetadataMsg = 'general.shortLoadingErrorMsg';
        this.longErrorMsg = 'general.longLoadingErrorMsg';
      }
    }));
  }

  private loadSupplyData(checkSavedData: boolean) {
    let oldSavedDataUsed = false;

    // Get the coin supply data.
    // Use saved data or get from the server. If there is no saved data, savedData is null.
    const savedData = checkSavedData ? this.getLocalValue(this.persistentServerSupplyResponseKey) : null;
    let nextOperation =  this.api.getCoinSupply();
    if (savedData) {
      nextOperation = of(savedData.value);
      oldSavedDataUsed = savedData.date < (new Date()).getTime() - dataValidityTime;
    }

    this.pageSubscriptions.push(nextOperation.pipe(first()).subscribe(response => {
      if (!savedData) {
        this.saveLocalValue(this.persistentServerSupplyResponseKey, response);
      }

      this.currentSupply = response.current_supply;
      this.totalSupply = response.total_supply;
      this.currentCoinhourSupply = response.current_coinhour_supply;
      this.totalCoinhourSupply = response.total_coinhour_supply;

      // If old saved data was used, repeat the operation, ignoring the saved data.
      if (oldSavedDataUsed) {
        this.loadSupplyData(false);
      }
    }, () => {
      // Error loading the data.
      this.loadingCoinSupplyMsg = 'general.shortLoadingErrorMsg';
    }));
  }

  ngOnDestroy() {
    // Clean all subscriptions.
    this.pageSubscriptions.forEach(sub => sub.unsubscribe());
    this.nodeUrlSubscription.unsubscribe();
  }
}
