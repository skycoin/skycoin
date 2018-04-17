import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { Wallet } from '../../../../app.datatypes';

@Component({
  selector: 'app-backup',
  templateUrl: './backup.component.html',
  styleUrls: ['./backup.component.css']
})
export class BackupComponent implements OnDestroy, OnInit {

  folder: string;

  constructor(
    public walletService: WalletService,
  ) {}

  ngOnInit() {
    this.walletService.folder().subscribe(folder => this.folder = folder);
  }

  ngOnDestroy() {
    this.walletService.all().subscribe(wallets => wallets.forEach(wallet => wallet.visible = false));
  }

  // download(wallet: Wallet) {
  //   const blob: Blob = new Blob([JSON.stringify({ seed: wallet.seed })], { type: 'application/json'});
  //   const link = document.createElement('a');
  //   link.href = window.URL.createObjectURL(blob);
  //   link['download'] = wallet.filename + '.json';
  //   link.click();
  // }
}
