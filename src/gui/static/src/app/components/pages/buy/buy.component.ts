/*
  IMPORTANT: Unused for a long time, it may need changes to work properly.
*/
import { filter, first } from 'rxjs/operators';
import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../services/wallet.service';
import { Address, PurchaseOrder, Wallet } from '../../../app.datatypes';
import { ButtonComponent } from '../../layout/button/button.component';
import { SubscriptionLike } from 'rxjs';
import { MsgBarService } from '../../../services/msg-bar.service';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.scss'],
})
export class BuyComponent implements OnInit, OnDestroy {
  @ViewChild('button', { static: false }) button: ButtonComponent;

  address: Address;
  config: any;
  form: FormGroup;
  order: PurchaseOrder;
  wallets: Wallet[];

  private subscriptionsGroup: SubscriptionLike[] = [];

  constructor(
    private formBuilder: FormBuilder,
    private purchaseService: PurchaseService,
    private msgBarService: MsgBarService,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
    this.loadData();
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
  }

  checkStatus() {
    this.button.setLoading();
    this.purchaseService.scan(this.order.recipient_address).subscribe(
      response => {
        this.button.setSuccess();
        this.order.status = response.status;
      },
      error => this.button.setError(error),
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

    this.subscriptionsGroup.push(this.form.get('wallet').valueChanges.subscribe(filename => {
      const wallet = this.wallets.find(wlt => wlt.filename === filename);
      console.log('changing wallet value', filename);
      this.purchaseService.generate(wallet).subscribe(
        order => this.saveData(order),
        error => this.msgBarService.showError(error.toString()),
      );
    }));
  }

  private loadConfig() {
    this.purchaseService.config().pipe(
      filter(config => !!config && !!config.sky_btc_exchange_rate), first())
      .subscribe(config => this.config = config);
  }

  private loadData() {
    this.loadConfig();
    this.loadOrder();

    this.subscriptionsGroup.push(this.walletService.all().subscribe(wallets => {
      this.wallets = wallets;

      if (this.order) {
        this.form.get('wallet').setValue(this.order.filename, { emitEvent: false });
      }
    }));
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
    this.purchaseService.scan(this.order.recipient_address).pipe(first()).subscribe(
      response => this.order.status = response.status,
      error => console.log(error),
    );
  }
}
