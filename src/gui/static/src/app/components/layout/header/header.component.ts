import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  @Input() title: string;
  @Input() coins: number;
  @Input() hours: number;

  get showLargeHeader(): boolean {
    return this.coins >= 0;
  }
}
