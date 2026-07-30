import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, filter, first } from 'rxjs';

// Add vars and functions to this file when having problems with circular dependencies
@Injectable()
export class GlobalsService {
  private nodeVersion: BehaviorSubject<string | null> = new BehaviorSubject<string | null>(null);

  setNodeVersion(version: string | null) {
    this.nodeVersion.next(version);
  }

  getValidNodeVersion(): Observable<string> {
    return this.nodeVersion.pipe(filter((version): version is string => version !== null), first());
  }
}
