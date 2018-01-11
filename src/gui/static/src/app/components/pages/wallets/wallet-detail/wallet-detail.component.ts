import { Component, Input } from '@angular/core';
import { Wallet } from "../../../../app.datatypes";

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss']
})
export class WalletDetailComponent {
  @Input() wallet: Wallet;
}
