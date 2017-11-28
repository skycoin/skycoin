import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  @Input() title: string;
  @Input() coins: number;
  @Input() hours: number;

  get showLargeHeader(): boolean {
    return this.coins >= 0;
  }

  constructor() { }

  ngOnInit() {
  }

}
