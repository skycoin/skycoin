import { Component, Input } from '@angular/core';

/**
 * Shows a loading animation. When the animation is not active, it shows an info msg and an
 * info icon, useful for indicating that no data was found or any other problem related to
 * loading the data.
 */
@Component({
  selector: 'app-loading-content',
  templateUrl: './loading-content.component.html',
  styleUrls: ['./loading-content.component.scss'],
})
export class LoadingContentComponent {
  // When true, the loading animation and a predefined loading msg is shown.
  @Input() isLoading = true;
  // Msg shown if isLoading is false.
  @Input() noDataText: string;
}
