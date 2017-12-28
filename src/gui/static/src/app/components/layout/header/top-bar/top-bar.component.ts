import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-top-bar',
  templateUrl: './top-bar.component.html',
  styleUrls: ['./top-bar.component.scss']
})
export class TopBarComponent {
  @Input() title: string;

  get showBack() {
    return this.router.url !== '/wallets';
  }

  constructor(
    private router: Router,
  ) {}

  goBack() {
    this.router.navigate(['/wallets']);
  }
}
