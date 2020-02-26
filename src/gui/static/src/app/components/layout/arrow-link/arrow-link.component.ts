import { Component, Input, Output, EventEmitter } from '@angular/core';

/**
 * Shows a link-like text with an arrow at the right. Used for showing more options or a list.
 */
@Component({
  selector: 'app-arrow-link',
  templateUrl: 'arrow-link.component.html',
  styleUrls: ['arrow-link.component.scss'],
})
export class ArrowLinkComponent {
  // Removes the padding at the left.
  @Input() noPadding = false;
  // Makes the arrow at the right to point up (false) or down (true).
  @Input() pointDown = true;
  @Input() text = '';
  @Output() pressed = new EventEmitter<any>();

  onClick(event) {
    this.pressed.emit(event);
  }
}
