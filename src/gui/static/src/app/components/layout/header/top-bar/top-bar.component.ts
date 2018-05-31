import { Component, Input } from '@angular/core';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-top-bar',
  templateUrl: './top-bar.component.html',
  styleUrls: ['./top-bar.component.scss'],
})
export class TopBarComponent {
  @Input() headline: string;

  constructor(
    public appService: AppService,
  ) {}
}
