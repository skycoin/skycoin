import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { AppService } from './services/app.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(
    private appService: AppService,
  ) {}

  ngOnInit() {
    this.appService.testBackend();
  }
}
