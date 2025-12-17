import { Component, EventEmitter, Input, OnInit, Output, ViewEncapsulation } from '@angular/core';

export enum DoubleButtonActive {RightButton, LeftButton}

@Component({
  selector: 'app-double-button',
  templateUrl: './double-button.component.html',
  styleUrls: ['./double-button.component.scss'],
  encapsulation: ViewEncapsulation.Emulated,
})
export class DoubleButtonComponent {
  @Input() rightButtonText: string;
  @Input() leftButtonText: string;
  @Input() activeButton: DoubleButtonActive;
  @Input() className = '';
  @Output() onStateChange = new EventEmitter();
  doubleButtonActive = DoubleButtonActive;

  onClick(button: DoubleButtonActive) {
    if (this.activeButton !== button) {
      this.activeButton = button;
      this.onStateChange.emit(button);
    }
  }
}
