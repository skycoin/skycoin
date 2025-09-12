import { throwError as observableThrowError, of as observableOf,  Observable } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';

import { ExplorerService } from '../explorer/explorer.service';

/**
 * Response returned by SearchService.processTerm. It has 2 properties, but only one will be
 * set, So it is necessary to check which one has a value.
 */
export class ResultNavCommandsResponse {
  /**
   * Observable for getting the navigation commands needed by Router.navigate for sending the
   * user to the page with the requested data. This Observable may fail due to different reasons.
   */
  resultNavCommands: Observable<string[]>;
  /**
   * Error returned if the function was not able to process the search term.
   */
  error: SearchError;
}

/**
 * List of errors that SearchService.processTerm may return.
 */
export enum SearchError {
  // The term searched is not valid.
  invalidSearchTerm = 1,
}

/**
 * Allows to process search terms to know where the user must be redirected to.
 */
@Injectable()
export class SearchService {

  constructor(
    private explorer: ExplorerService,
  ) { }

  /**
   * This function is used to process a search term entered by the user, to get the navigation commands
   * needed by Router.navigate for sending the user to the page with the requested data.
   *
   * @param searchTerm search term to be processed.
   * @returns An error or an observable that will get get the navigation commands needed by Router.navigate
   * for sending the user to the page with the requested data.
   */
  processTerm(searchTerm: string): ResultNavCommandsResponse {
    const response = new ResultNavCommandsResponse();
    searchTerm = encodeURIComponent(searchTerm);

    if (searchTerm.length >= 27 && searchTerm.length <= 35) {
      // If the search term has this length, it is assumed to be an address.
      response.resultNavCommands = observableOf(['/app/address', searchTerm]);
    } else if (searchTerm.length === 64) {
      // If the search term has this length, it is a hash, but we don't know if it is the hash of
      // a block or a transaction, so we try to find a block with that hash.
      response.resultNavCommands = this.explorer.getBlockByHash(searchTerm).pipe(
        // If a block with that hash is found, it is assumed that the user was searching for that block.
        map(block => ['/app/block', block.id.toString()]),
        catchError((error: any) => {
          // If the node returns 404 (no block with an ID equal to the searched hash was found) then
          // it is assumed that the user was searching for a transaction.
          if (error && error.status && error.status === 404) {
            return observableOf(['/app/transaction', searchTerm]);
          } else {
            // If any other error happened, the situation is unexpected and the Observable fails.
            return observableThrowError(error);
          }
        }));
    } else {
      // If the search term is a positive integer, it is assumed to be a block number.
      if (parseInt(searchTerm, 10).toString() === searchTerm && parseInt(searchTerm, 10) >= 0) {
        response.resultNavCommands = observableOf(['/app/block', searchTerm]);
      } else {
        // If the search term did not match any of the previous criteria, it is not possible to process it.
        response.error = SearchError.invalidSearchTerm;
      }
    }

    return response;
  }

}
