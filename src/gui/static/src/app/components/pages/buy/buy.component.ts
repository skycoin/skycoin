/*
  IMPORTANT: Unused for a long time, it may need changes to work properly.
*/
import { filter, first } from 'rxjs';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { PurchaseOrder } from '../../../app.datatypes';
import { ButtonComponent } from '../../layout/button/button.component';
import { SubscriptionLike } from 'rxjs';
import { MsgBarService } from '../../../services/msg-bar.service';
import { WalletBase, AddressBase } from '../../../services/wallet-operations/wallet-objects';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';

@Component({
    selector: 'app-buy',
    templateUrl: './buy.component.html',
    styleUrls: ['./buy.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class BuyComponent implements OnInit, OnDestroy {
  @ViewChild('button') button!: ButtonComponent;

  address!: AddressBase;
  config: any;
  form!: UntypedFormGroup;
  order: PurchaseOrder | null = null;
  wallets!: WalletBase[];

  private subscriptionsGroup: SubscriptionLike[] = [];

  constructor(
    private formBuilder: UntypedFormBuilder,
    private purchaseService: PurchaseService,
    private msgBarService: MsgBarService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private changeDetectorRef: ChangeDetectorRef,
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
    this.purchaseService.scan(this.order!.recipient_address).subscribe(
      response => {
        this.button.setSuccess();
        this.order!.status = response.status;
        this.changeDetectorRef.markForCheck();
      },
      // On this part the error was shown on the button. Now it would have to be shown on the msg bar.
      error => {
        this.button.resetState();
        this.changeDetectorRef.markForCheck();
      },
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

    this.subscriptionsGroup.push(this.form.get('wallet')!.valueChanges.subscribe(id => {
      const wallet = this.wallets.find(wlt => wlt.id === id);
      console.log('changing wallet value', id);
      this.purchaseService.generate(wallet!).subscribe(
        order => {
          this.saveData(order);
          this.changeDetectorRef.markForCheck();
        },
        error => {
          this.msgBarService.showError(error.toString());
          this.changeDetectorRef.markForCheck();
        },
      );
      this.changeDetectorRef.markForCheck();
    }));
  }

  private loadConfig() {
    this.purchaseService.config().pipe(
      filter(config => !!config && !!config.sky_btc_exchange_rate), first())
      .subscribe(config => {
        this.config = config;
        this.changeDetectorRef.markForCheck();
      });
  }

  private loadData() {
    this.loadConfig();
    this.loadOrder();

    this.subscriptionsGroup.push(this.walletsAndAddressesService.allWallets.subscribe(wallets => {
      this.wallets = wallets;

      if (this.order) {
        this.form.get('wallet')!.setValue(this.order.filename, { emitEvent: false });
      }
      this.changeDetectorRef.markForCheck();
    }));
  }

  private loadOrder() {
    const order: PurchaseOrder = JSON.parse(window.localStorage.getItem('purchaseOrder')!);
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
    this.purchaseService.scan(this.order!.recipient_address).pipe(first()).subscribe(
      response => {
        this.order!.status = response.status;
        this.changeDetectorRef.markForCheck();
      },
      error => {
        console.log(error);
        this.changeDetectorRef.markForCheck();
      },
    );
  }
}
