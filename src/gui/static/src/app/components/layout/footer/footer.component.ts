import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

  openWalletsPage() {
    // this.router(WalletsPage);
  }

  openSendPage() {
    // const modal = this.modal.create(SendSkycoinPage);
    // modal.present();
  }

  openTransactions() {
    // this.router(TransactionsPage);
  }
}
