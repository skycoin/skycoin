import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { AppService } from './services/app.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(
    private appService: AppService,
    private translateService: TranslateService,
  ) {
    translateService.setDefaultLang('en');
    translateService.use('en');
  }

  ngOnInit() {
    this.appService.testBackend();
  }
}
