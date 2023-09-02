import { Component, EventEmitter, Input, Output, ViewEncapsulation, OnDestroy } from '@angular/core';

/**
 * Identifies the active button of a DoubleButtonComponent.
 */
export enum DoubleButtonActive {
  RightButton = 'RightButton',
  LeftButton = 'LeftButton',
}

/**
 * Bar with 2 buttons, one active and other inactive. Used to select between two options.
 */
@Component({
  selector: 'app-double-button',
  templateUrl: './double-button.component.html',
  styleUrls: ['./double-button.component.scss'],
  encapsulation: ViewEncapsulation.Emulated,
})
export class DoubleButtonComponent implements OnDestroy {
  @Input() rightButtonText: string;
  @Input() leftButtonText: string;
  @Input() activeButton: DoubleButtonActive;
  // Allows to add classes to the component.
  @Input() className = '';
  // If true, when the user clicks one of the buttons the newly selected button will not
  // be selected automatically. Instead, the control will just send the event indicating the
  // click and the "activeButton" property will have to be changed for the clicked button to
  // be selected.
  @Input() changeActiveButtonManually = false;
  @Output() stateChange = new EventEmitter<DoubleButtonActive>();
  ButtonStates = DoubleButtonActive;

  ngOnDestroy() {
    this.stateChange.complete();
  }

  onRightButtonClicked() {
    if (this.activeButton === DoubleButtonActive.LeftButton) {
      if (!this.changeActiveButtonManually) {
        this.activeButton = DoubleButtonActive.RightButton;
      }
      this.stateChange.emit(DoubleButtonActive.RightButton);
    }
  }

  onLeftButtonClicked() {
    if (this.activeButton === DoubleButtonActive.RightButton) {
      if (!this.changeActiveButtonManually) {
        this.activeButton = DoubleButtonActive.LeftButton;
      }
      this.stateChange.emit(DoubleButtonActive.LeftButton);
    }
  }
}
