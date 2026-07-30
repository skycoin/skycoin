import { Component, ChangeDetectionStrategy } from '@angular/core';

@Component({
    selector: 'app-disclaimer-warning',
    templateUrl: './disclaimer-warning.component.html',
    styleUrls: ['./disclaimer-warning.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class DisclaimerWarningComponent {

  constructor() { }
}
