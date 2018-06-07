import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class NavBarService {
  switchVisible = false;
  activeComponent = new BehaviorSubject(1);

  setActiveComponent(value) {
    this.activeComponent.next(value);
  }

  showSwitch() {
    this.switchVisible = true;
  }

  hideSwitch() {
    this.switchVisible = false;
  }
}
