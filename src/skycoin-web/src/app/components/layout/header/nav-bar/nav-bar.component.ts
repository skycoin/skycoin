import { Component, ChangeDetectionStrategy } from '@angular/core';

import { NavBarService } from '../../../../services/nav-bar.service';

@Component({
    selector: 'app-nav-bar',
    templateUrl: './nav-bar.component.html',
    styleUrls: ['./nav-bar.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class NavBarComponent {
  constructor(
    public navbarService: NavBarService,
  ) { }

  changeActiveComponent(value: any) {
    this.navbarService.setActiveComponent(value);
  }
}
