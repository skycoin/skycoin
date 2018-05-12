import { Component, EventEmitter, Input, Output, ViewEncapsulation } from '@angular/core';

export enum DoubleButtonActive { RightButton, LeftButton }

@Component({
  selector: 'app-double-button',
  templateUrl: './double-button.component.html',
  styleUrls: ['./double-button.component.scss'],
  encapsulation: ViewEncapsulation.Emulated,
})
export class DoubleButtonComponent {
  @Input() rightButtonText: any;
  @Input() leftButtonText: any;
  @Input() activeButton: DoubleButtonActive;
  @Output() onStateChange = new EventEmitter();
  ButtonState = DoubleButtonActive;

  onRightClick() {
    if (this.activeButton === DoubleButtonActive.LeftButton) {
      this.activeButton = DoubleButtonActive.RightButton;
      this.onStateChange.emit(DoubleButtonActive.RightButton);
    }
  }

  onLeftClick() {
    if (this.activeButton === DoubleButtonActive.RightButton) {
      this.activeButton = DoubleButtonActive.LeftButton;
      this.onStateChange.emit(DoubleButtonActive.LeftButton);
    }
  }
}
