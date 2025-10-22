import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable()
export class NavBarService {
  switchVisible = false;
  activeComponent = new BehaviorSubject(1);
  leftText: string;
  rightText: string;

  setActiveComponent(value = 1) {
    this.activeComponent.next(value);
  }

  showSwitch(leftText: string, rightText: string, initialSwitchState: number | null = null) {
    this.switchVisible = true;
    this.leftText = leftText;
    this.rightText = rightText;

    if (initialSwitchState) {
      this.setActiveComponent(initialSwitchState);
    }
  }

  hideSwitch() {
    this.switchVisible = false;
  }
}
