import { Component } from '@angular/core';
import { MatDialogRef, MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { WalletService } from '../../../../services/wallet.service';
import { Wallet, ConfirmationData } from '../../../../app.datatypes';
import { ApiService } from '../../../../services/api.service';
import { showConfirmationModal } from '../../../../utils';

@Component({
  selector: 'app-show-wallet',
  templateUrl: './show-wallet.component.html',
  styleUrls: ['./show-wallet.component.scss'],
})
export class ShowWalletComponent {

  wallets: Wallet[] = [];

  constructor(
    public dialogRef: MatDialogRef<ShowWalletComponent>,
    private walletService: WalletService,
    private translateService: TranslateService,
    private dialog: MatDialog,
    apiService: ApiService,
  ) {
    apiService.getWallets().first().subscribe(recoveredWallets =>
      this.wallets = recoveredWallets.filter(wallet => walletService.hiddenWalletsMap.has(wallet.filename)),
    );
  }

  unhide(wallet: Wallet) {
    const confirmationData: ConfirmationData = {
      text: this.translateService.instant('wallet.show-wallet.show-confirmation', {name: wallet.label}),
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.dialogRef.close(this.walletService.unhideWallet(wallet.filename));
      }
    });
  }
}
