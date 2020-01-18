import { Component } from '@angular/core';
import { AppConfig } from '../../../../app.config';
import { NavBarService } from '../../../../services/nav-bar.service';
import { environment } from '../../../../../environments/environment';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.scss'],
})
export class NavBarComponent {
  otcEnabled = AppConfig.otcEnabled;
  exchangeEnabled = !!environment.swaplab.apiKey;

  constructor(
    public appService: AppService,
    public navbarService: NavBarService,
  ) { }

  changeActiveComponent(value) {
    this.navbarService.setActiveComponent(value);
  }
}
