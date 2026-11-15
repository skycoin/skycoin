import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, filter, first, throwError, timeout } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

// Add vars and functions to this file when having problems with circular dependencies
@Injectable()
export class GlobalsService {
  /** How long an operation waits for the first check of the node to pass. */
  nodeVersionWaitMs = 30000;

  private nodeVersion: BehaviorSubject<string | null> = new BehaviorSubject<string | null>(null);

  constructor(private translate: TranslateService) {}

  setNodeVersion(version: string | null) {
    this.nodeVersion.next(version);
  }

  /** The node version. Fails when it is not known in time, rather than waiting forever. */
  getValidNodeVersion(): Observable<string> {
    return this.nodeVersion.pipe(
      filter((version): version is string => version !== null),
      first(),
      timeout({
        first: this.nodeVersionWaitMs,
        with: () => throwError(() => new Error(this.translate.instant('service.api.node-not-ready'))),
      }),
    );
  }
}
