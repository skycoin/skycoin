import { Component } from '@angular/core';

import { FooterConfig } from 'app/app.config';

/**
 * Generic customizable footer. To activate it, FooterConfig.useGenericFooter (in app.config.ts)
 * must be true. Read the docs for more information.
 */
@Component({
  selector: 'app-generic-footer',
  templateUrl: './generic-footer.component.html',
  styleUrls: ['./generic-footer.component.scss']
})
export class GenericFooterComponent {
  /**
   * Configuration params.
   */
  config = FooterConfig;
}
