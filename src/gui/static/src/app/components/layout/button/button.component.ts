import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-button',
  templateUrl: 'button.component.html',
  styleUrls: ['button.component.scss']
})

export class ButtonComponent {
  @Input() disabled: any;
  @Output() action = new EventEmitter();

  error: string;
  state: number;

  onClick() {
    if (!this.disabled) {
      this.action.emit();
    }
  }

  setLoading() {
    this.state = 0;
  }

  setSuccess() {
    this.state = 1;
    setTimeout(() => this.state = null, 3000);
  }

  setError(error: any) {
    this.error = typeof error === 'string' ? error : error['_body'];
    this.state = 2;
  }

  resetState() {
    this.state = null;
    this.error = '';
  }
}
