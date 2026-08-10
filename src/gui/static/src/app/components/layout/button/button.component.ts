import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';

enum ButtonStates {
  Normal = 'Normal',
  Loading = 'Loading',
  Success = 'Success',
}

/**
 * Normal rounded button used in most parts of the app.
 */
@Component({
    selector: 'app-button',
    templateUrl: 'button.component.html',
    styleUrls: ['button.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class ButtonComponent implements OnDestroy {
  constructor(private changeDetectorRef: ChangeDetectorRef) {
  }

  @Input() disabled!: boolean;
  // If true, the button will send click events even when disabled.
  @Input() forceEmitEvents = false;
  // Click event.
  @Output() action = new EventEmitter();
  @ViewChild('button') button!: HTMLButtonElement;

  state = ButtonStates.Normal;
  buttonStates = ButtonStates;

  ngOnDestroy() {
    this.action.complete();
  }

  onClick() {
    if (!this.disabled || this.forceEmitEvents) {
      this.action.emit();
    }
  }

  /**
   * Focuses the button.
   */
  focus() {
    this.button.focus();
  }

  /**
   * Shows the loading animation. The button does not send click events while the
   * animation is active.
   */
  setLoading() {
    this.state = ButtonStates.Loading;
  }

  /**
   * Shows the success icon.
   */
  setSuccess() {
    this.state = ButtonStates.Success;
    setTimeout(() => {
      this.state = ButtonStates.Normal;
      this.changeDetectorRef.markForCheck();
    }, 3000);
  }

  setDisabled() {
    this.disabled = true;
  }

  setEnabled() {
    this.disabled = false;
  }

  isLoading(): boolean {
    return this.state === ButtonStates.Loading;
  }

  /**
   * Removes the icons and animations, but does not affects the enabled/disabled status.
   * @returns The currents instance is returned to make it easier to concatenate function calls.
   */
  resetState() {
    this.state = ButtonStates.Normal;

    return this;
  }
}
