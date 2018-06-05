import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class NavBarService {
  switchVisible = false;
  activeComponent = new BehaviorSubject(1);
  leftText: string;
  rightText: string;

  setActiveComponent(value) {
    this.activeComponent.next(value);
  }

  showSwitch(leftText, rightText) {
    this.switchVisible = true;
    this.leftText = leftText;
    this.rightText = rightText;
  }

  hideSwitch() {
    this.switchVisible = false;
  }
}
