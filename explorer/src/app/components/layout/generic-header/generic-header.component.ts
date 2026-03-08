import { Component } from '@angular/core';

import { HeaderConfig } from 'app/app.config';

/**
 * Generic customizable header. To activate it, HeaderConfig.useGenericHeader (in app.config.ts)
 * must be true. Read the docs for more information.
 */
@Component({
    selector: 'app-generic-header',
    templateUrl: './generic-header.component.html',
    styleUrls: ['./generic-header.component.scss'],
    standalone: false
})
export class GenericHeaderComponent {
  /**
   * Configuration params.
   */
  config = HeaderConfig;
}
