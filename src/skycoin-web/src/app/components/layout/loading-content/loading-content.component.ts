import { Component, Input, ChangeDetectionStrategy } from '@angular/core';

@Component({
    selector: 'app-loading-content',
    templateUrl: './loading-content.component.html',
    styleUrls: ['./loading-content.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class LoadingContentComponent {
  @Input() isLoading = true;
  @Input() showError = false;
  @Input() noDataText = 'tx.none';
  @Input() errorText = 'errors.loading-error';
}
