import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {

  constructor(
    private router: Router,
  ) { }

  openWalletsPage() {
    this.router.navigate(['/wallets']);
  }

  openSendPage() {
    // const modal = this.modal.create(SendSkycoinPage);
    // modal.present();
  }

  openTransactions() {
    // this.router(TransactionsPage);
  }
}
