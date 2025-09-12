import { Component, Input } from '@angular/core';

/**
 * Takes a date string in an Unix like format and displays it in an easy-to-read format, clearly
 * separating the date from the hour and showing timezone info.
 */
@Component({
  selector: 'app-date-formatter',
  templateUrl: './date-formatter.component.html',
  styleUrls: ['./date-formatter.component.scss']
})
export class DateFormatterComponent {
  /**
   * Date to be displayed. It should expresed in time since the epoch.
   */
  @Input()date;

  /**
   * Number with which the value of "date" must be multiplied to convert it into milliseconds since
   * the epoch. For example, if "date" has a value expresed in seconds since the epoch, the value
   * must be multiplied by 1000 to conver it into milliseconds since the epoch.
   */
  @Input()dateMultiplier = 1000;

  /**
   * Indicates whether the mouse cursor is over the control or not.
   */
  mouseOver = false;

  /**
   * Event called when the mouse cursor enters the control. The event is associated in the
   * HTML file.
   */
  onMouseEnter() {
    this.mouseOver = true;
  }

  /**
   * Event called when the mouse cursor leaves the control. The event is associated in the
   * HTML file.
   */
  onMouseLeave() {
    this.mouseOver = false;
  }
}
