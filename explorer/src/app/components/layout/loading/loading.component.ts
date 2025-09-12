import { Component, Input } from '@angular/core';

/**
 * Control for showing a loading animation. It is also used to show an error message if the
 * operation fails.
 */
@Component({
  selector: 'app-loading',
  templateUrl: './loading.component.html',
  styleUrls: ['./loading.component.scss']
})
export class LoadingComponent {
  /**
   * Message to show with the loading animation.
   */
  loadingMsg = 'general.waitingData';
  /**
   * Message to show in case of error. If set, the loading animation is removed.
   */
  @Input() longErrorMsg: string;
}
