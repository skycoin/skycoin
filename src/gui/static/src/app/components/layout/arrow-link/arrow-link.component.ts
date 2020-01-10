import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-arrow-link',
  templateUrl: 'arrow-link.component.html',
  styleUrls: ['arrow-link.component.scss'],
})
export class ArrowLinkComponent {
  @Input() noPadding = false;
  @Input() pointDown = true;
  @Input() text = '';
  @Output() pressed = new EventEmitter<any>();

  onClick(event) {
    this.pressed.emit(event);
  }
}
