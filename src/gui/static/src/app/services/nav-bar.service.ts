import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { DoubleButtonActive } from '../components/layout/double-button/double-button.component';

@Injectable()
export class NavBarService {
  switchVisible = false;
  activeComponent = new BehaviorSubject(1);
  leftText: string;
  rightText: string;
  switchDiabled = false;

  setActiveComponent(value) {
    this.activeComponent.next(value);
  }

  showSwitch(leftText, rightText, selectedButton = DoubleButtonActive.LeftButton) {
    this.setActiveComponent(selectedButton);
    this.switchDiabled = false;
    this.switchVisible = true;
    this.leftText = leftText;
    this.rightText = rightText;
  }

  hideSwitch() {
    this.switchVisible = false;
  }

  enableSwitch() {
    this.switchDiabled = false;
  }

  disableSwitch() {
    this.switchDiabled = true;
  }
}
