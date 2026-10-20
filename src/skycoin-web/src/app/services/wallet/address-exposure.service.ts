import { Injectable } from '@angular/core';
import { Observable, from, of } from 'rxjs';
import { catchError, map, mergeMap, toArray } from 'rxjs';

import { ApiService } from '../api.service';

/**
 * Finds addresses that have sent coins. A spend puts the address's public key on the
 * blockchain, so a future quantum computer could take what is still sent to it.
 */
@Injectable()
export class AddressExposureService {
  // A spent address stays spent, so it is cached for good. Unspent ones are rechecked
  // after unspentTtl.
  private spent = new Set<string>();
  private unspentCheckedAt = new Map<string, number>();
  private readonly unspentTtl = 5 * 60 * 1000;

  constructor(private apiService: ApiService) { }

  /**
   * Returns the addresses in the list that have sent coins, asking the current coin's
   * node. Addresses the node cannot answer for are left out.
   */
  getExposed(addresses: string[]): Observable<Set<string>> {
    const now = Date.now();
    const toCheck = Array.from(new Set(addresses)).filter(a =>
      !this.spent.has(a) && now - (this.unspentCheckedAt.get(a) || 0) > this.unspentTtl,
    );

    return from(toCheck).pipe(
      mergeMap(address => this.apiService.get('address_uxouts', { address }).pipe(
        map((uxouts: any[]) => {
          if ((uxouts || []).some(ux => ux.spent_block_seq > 0)) {
            this.spent.add(address);
          } else {
            this.unspentCheckedAt.set(address, Date.now());
          }
        }),
        catchError(() => of(null)),
      ), 4),
      toArray(),
      map(() => new Set(addresses.filter(a => this.spent.has(a)))),
    );
  }
}
