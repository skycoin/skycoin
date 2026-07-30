import { Component, ChangeDetectionStrategy } from '@angular/core';

import { AppConfig } from '../../../../app.config';
import { NavBarSwitchService } from '../../../../services/nav-bar-switch.service';
import { environment } from '../../../../../environments/environment';
import { AppService } from '../../../../services/app.service';

/**
 * Navigation bar shown on the header.
 */
@Component({
    selector: 'app-nav-bar',
    templateUrl: './nav-bar.component.html',
    styleUrls: ['./nav-bar.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class NavBarComponent {
  otcEnabled = AppConfig.otcEnabled;
  exchangeEnabled = !!environment.swaplab.apiKey;

  constructor(
    public appService: AppService,
    public navBarSwitchService: NavBarSwitchService,
  ) { }

  changeActiveComponent(value: any) {
    this.navBarSwitchService.setActiveComponent(value);
  }
}
