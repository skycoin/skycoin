import { Component, ChangeDetectionStrategy } from '@angular/core';

@Component({
    selector: 'app-disclaimer-warning',
    templateUrl: './disclaimer-warning.component.html',
    styleUrls: ['./disclaimer-warning.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class DisclaimerWarningComponent {

  constructor() { }
}
