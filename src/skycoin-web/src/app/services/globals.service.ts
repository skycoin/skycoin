import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { BehaviorSubject } from 'rxjs';

// Add vars and functions to this file when having problems with circular dependencies
@Injectable()
export class GlobalsService {
  private nodeVersion: BehaviorSubject<string> = new BehaviorSubject<string>(null);

  setNodeVersion(version: string) {
    this.nodeVersion.next(version);
  }

  getValidNodeVersion(): Observable<string> {
    return this.nodeVersion.filter(version => version !== null).first();
  }
}
