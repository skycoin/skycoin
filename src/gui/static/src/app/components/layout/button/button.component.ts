import { Component, EventEmitter, Input, Output, ViewChild, ElementRef } from '@angular/core';
import { parseResponseMessage } from '../../../utils/errors';
import { MatTooltip } from '@angular/material';

@Component({
  selector: 'app-button',
  templateUrl: 'button.component.html',
  styleUrls: ['button.component.scss'],
})
export class ButtonComponent {
  @Input() disabled: any;
  @Input() forceEmitEvents = false;
  @Output() action = new EventEmitter();
  @ViewChild('tooltip') tooltip: MatTooltip;
  @ViewChild('button') button: HTMLButtonElement;

  error: string;
  state: number;
  mouseOver = false;

  onClick() {
    if (!this.disabled || this.forceEmitEvents) {
      this.error = '';
      this.action.emit();
    }
  }

  focus() {
    this.button.focus();
  }

  setLoading() {
    this.state = 0;
  }

  setSuccess() {
    this.state = 1;
    setTimeout(() => this.state = null, 3000);
  }

  setError(error: any) {
    this.error = typeof error === 'string' ? error : parseResponseMessage(error['_body']);
    this.state = 2;

    if (this.mouseOver) {
      setTimeout(() => this.tooltip.show(), 50);
    }
  }

  setDisabled() {
    this.disabled = true;
  }

  setEnabled() {
    this.disabled = false;
  }

  isLoading() {
    return this.state === 0;
  }

  resetState() {
    this.state = null;
    this.error = '';

    return this;
  }
}
