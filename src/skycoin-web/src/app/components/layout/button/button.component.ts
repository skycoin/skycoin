import { Component, EventEmitter, Input, Output, ViewChild, ChangeDetectionStrategy } from '@angular/core';
import { MatTooltip } from '@angular/material/tooltip';

@Component({
    selector: 'app-button',
    templateUrl: 'button.component.html',
    styleUrls: ['button.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})

export class ButtonComponent {
  @Input() disabled!: boolean;
  @Input() forceEmitEvents = false;
  @Input() spinnerStyle = 'primary';
  @Output() action = new EventEmitter();
  @ViewChild('tooltip') tooltip!: MatTooltip;

  error!: string;
  state!: number | null;
  mouseOver = false;

  onClick() {
    if (!this.disabled || this.forceEmitEvents) {
      this.error = '';
      this.action.emit();
    }
  }

  setLoading() {
    this.state = 0;
  }

  isLoading(): boolean {
    return this.state === 0;
  }

  setSuccess() {
    this.state = 1;
    setTimeout(() => this.state = null, 3000);
  }

  setError(error: string) {
    this.error = error;
    this.state = 2;

    setTimeout(() => {
      if (this.mouseOver) {
        this.tooltip.show(50);
      }
    }, 0);
  }

  setEnabled() {
    this.disabled = false;
  }

  setDisabled() {
    this.disabled = true;
  }

  resetState() {
    this.state = null;
    this.error = '';
  }
}
