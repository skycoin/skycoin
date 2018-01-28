import { Component, ViewChild } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../services/wallet.service';
import { Address, PurchaseOrder, Wallet } from '../../../app.datatypes';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ButtonComponent } from '../../layout/button/button.component';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.scss']
})
export class BuyComponent {
  @ViewChild('button') button: ButtonComponent;

  address: Address;
  config: any;
  form: FormGroup;
  order: PurchaseOrder;
  wallets: Wallet[];

  constructor(
    private formBuilder: FormBuilder,
    private purchaseService: PurchaseService,
    private snackBar: MatSnackBar,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
    this.loadData();
  }

  checkStatus() {
    this.button.setLoading();
    this.purchaseService.scan(this.order.recipient_address).first().subscribe(
      response => {
        this.button.setSuccess();
        this.order.status = response.status;
      },
      error => this.button.setError(error)
    );
  }

  removeOrder() {
    window.localStorage.removeItem('purchaseOrder');
    this.order = null;
  }

  private initForm() {
    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
    });

    this.form.controls.wallet.valueChanges.subscribe(filename => {
      const wallet = this.wallets.find(wallet => wallet.filename === filename);
      console.log('changing wallet value', filename);
      this.purchaseService.generate(wallet).subscribe(
        order => this.saveData(order),
        error => this.snackBar.open(error.toString())
      );
    })
  }

  private loadConfig() {
    this.purchaseService.config()
      .filter(config => !!config && !!config.sky_btc_exchange_rate)
      .first()
      .subscribe(config => this.config = config);
  }

  private loadData() {
    this.loadConfig();
    this.loadOrder();

    this.walletService.all().subscribe(wallets => {
      this.wallets = wallets;

      if (this.order) {
        this.form.controls.wallet.setValue(this.order.filename, { emitEvent: false });
      }
    });
  }

  private loadOrder() {
    const order: PurchaseOrder = JSON.parse(window.localStorage.getItem('purchaseOrder'));
    if (order) {
      this.order = order;
      this.updateOrder();
    }
  }

  private saveData(order: PurchaseOrder) {
    this.order = order;
    window.localStorage.setItem('purchaseOrder', JSON.stringify(order));
  }

  private updateOrder() {
    this.purchaseService.scan(this.order.recipient_address).first().subscribe(
      response => this.order.status = response.status,
      error => console.log(error)
    );
  }
}
