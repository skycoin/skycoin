import { Component, Input } from '@angular/core';
import { trigger, state, style, transition, animate } from '@angular/animations';

/**
 * Button that allow the user to copy text.
 */
@Component({
  selector: 'app-copy-button',
  templateUrl: './copy-button.component.html',
  styleUrls: ['./copy-button.component.scss'],
  animations: [
    trigger('showMessage', [
      state('show', style({
        opacity: 1.0,
        transform: 'translateY(-25px)'
      })),
      state('hide', style({
        opacity: 0.0,
        transform: 'translateY(-35px)'
      })),
      state('reset', style({
        opacity: 0.0,
        transform: 'translateY(-15px)'
      })),
      transition('* => show', [
        animate('200ms 0ms ease-out')
      ]),
      transition('* => hide', [
        animate('200ms 500ms ease-in')
      ]),
    ])
  ]
})
export class CopyButtonComponent {

  /**
   * State for showing the black "Copied" message box.
   */
  private static showAnimState = 'show';
  /**
   * State for hidding the black "Copied" message box.
   */
  private static hideAnimState = 'hide';
  /**
   * Special state for immediately hidding the black "Copied" message box.
   */
  private static resetAnimState = 'reset';

  /**
   * Current animation state of the black "Copied" message box.
   */
  animState = CopyButtonComponent.resetAnimState;
  /**
   * If the black "Copied" message box should be visible.
   */
  showLabel = false;
  /**
   * Text to be copied when clicking this button.
   */
  @Input() text: string;

  /**
   * Function for copying the text.
   */
  copy() {
    // Create a temporary textarea for copying the text, add the text to it and place it in
    // a corner.
    const tempElement = document.createElement('textarea');
    tempElement.style.position = 'fixed';
    tempElement.style.left = '1px';
    tempElement.style.top = '1px';
    tempElement.style.width = '1px';
    tempElement.style.height = '1px';
    tempElement.style.opacity = '0';
    tempElement.value = this.text;
    document.body.appendChild(tempElement);

    // Select the text.
    tempElement.focus();
    tempElement.select();

    // Copy the text and remove the textarea.
    document.execCommand('copy');
    document.body.removeChild(tempElement);

    // Reset the animation of the black "Copied" message box and hide it.
    this.animState = CopyButtonComponent.resetAnimState;
    this.showLabel = false;

    // Wait a frame to start the animation for showing the black "Copied" message box.
    setTimeout(() => {
      this.animState = CopyButtonComponent.showAnimState;
      this.showLabel = true;
    }, 16);
  }

  /**
   * Event dispatched when an animation ends. The event is associated in the HTML file.
   */
  animationDone() {
    if (this.animState === CopyButtonComponent.showAnimState) {
      // Go to the next state.
      this.animState = CopyButtonComponent.hideAnimState;
    } else if (this.animState === CopyButtonComponent.hideAnimState) {
      // Hide the the black "Copied" message box.
      this.showLabel = false;
    }
  }
}
