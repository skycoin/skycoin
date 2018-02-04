import { Component, OnInit } from '@angular/core';
import { AppConfig } from '../../../../app.config';

@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.scss']
})
export class NavBarComponent implements OnInit {

  otcEnabled = AppConfig.otcEnabled;

  constructor() { }

  ngOnInit() {
  }

}
