import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

import { DoubleButtonActive } from '../components/layout/double-button/double-button.component';

/**
 * Allows to make operations related to the switch shown at the right side of the navigation bar.
 */
@Injectable()
export class NavBarSwitchService {

  // Properties
  ///////////////////////////////////////////

  /**
   * Indicates if the navigation bar must show the 2 options switch on the right side.
   */
  get switchVisible(): boolean {
    return this.switchVisibleInternal;
  }
  private switchVisibleInternal = false;

  /**
   * Indicates which option is selected.
   */
  get activeComponent(): Observable<DoubleButtonActive> {
    return this.activeComponentInternal.asObservable();
  }
  private activeComponentInternal = new BehaviorSubject(DoubleButtonActive.LeftButton);

  /**
   * Text of the left option.
   */
  get leftText(): string {
    return this.leftTextInternal;
  }
  private leftTextInternal: string;

  /**
   * Text of the right option.
   */
  get rightText(): string {
    return this.rightTextInternal;
  }
  private rightTextInternal: string;

  /**
   * Indicates if the switch must be shown disabled.
   */
  get switchDiabled(): boolean {
    return this.switchDiabledInternal;
  }
  private switchDiabledInternal = false;

  // functions
  ///////////////////////////////////////////

  /**
   * Indicates which option must be shown selected.
   */
  setActiveComponent(value: DoubleButtonActive) {
    this.activeComponentInternal.next(value);
  }

  /**
   * Makes the switch visible.
   * @param leftText Text for the option at the left.
   * @param rightText Text for the option at the right.
   * @param selectedButton Which option must be shown selected.
   */
  showSwitch(leftText, rightText, selectedButton = DoubleButtonActive.LeftButton) {
    this.setActiveComponent(selectedButton);
    this.switchDiabledInternal = false;
    this.switchVisibleInternal = true;
    this.leftTextInternal = leftText;
    this.rightTextInternal = rightText;
  }

  /**
   * Hides the switch.
   */
  hideSwitch() {
    this.switchVisibleInternal = false;
  }

  /**
   * Shows the switch enabled.
   */
  enableSwitch() {
    this.switchDiabledInternal = false;
  }

  /**
   * Shows the switch disabled.
   */
  disableSwitch() {
    this.switchDiabledInternal = true;
  }
}
