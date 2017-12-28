import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { SendSkycoinComponent } from '../../pages/send-skycoin/send-skycoin.component';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {

  constructor(
    private dialog: MdDialog,
    private router: Router,
  ) { }

  openWalletsPage() {
    this.router.navigate(['/wallets']);
  }

  openSendPage() {
    const config = new MdDialogConfig();
    config.width = '566px';
    this.dialog.open(SendSkycoinComponent, config);
  }

  openTransactions() {
    this.router.navigate(['/transactions']);
  }
}
