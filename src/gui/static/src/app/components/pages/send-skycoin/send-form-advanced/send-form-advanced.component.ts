import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog, MatSnackBar } from '@angular/material';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { showSnackbarError } from '../../../../utils/errors';
import { Subscription } from 'rxjs/Subscription';
import { NavBarService } from '../../../../services/nav-bar.service';

@Component({
  selector: 'app-send-form-advanced',
  templateUrl: './send-form-advanced.component.html',
  styleUrls: ['./send-form-advanced.component.scss'],
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('previewButton') previewButton: ButtonComponent;
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  form: FormGroup;
  addresses = [];
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx: boolean;

  private subscriptions: Subscription;

  constructor(
    public walletService: WalletService,
    private formBuilder: FormBuilder,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
    private navbarService: NavBarService,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced');

    this.form = this.formBuilder.group({
      wallet: ['', Validators.required],
      addresses: ['', Validators.required],
      changeAddress: [''],
      destinations: this.formBuilder.array(
        [this.createDestinationFormGroup()],
        this.validateDestinations.bind(this),
      ),
    });

    this.subscriptions = this.form.get('wallet').valueChanges.subscribe(wallet => {
      this.addresses = wallet.addresses.filter(addr => addr.coins > 0);
      this.form.get('addresses').setValue([]);
      this.form.get('destinations').updateValueAndValidity();
    });

    this.subscriptions.add(this.form.get('addresses').valueChanges.subscribe(() => {
      this.form.get('destinations').updateValueAndValidity();
    }));

    if (this.formData) {
      this.fillForm();
    }
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
    this.navbarService.hideSwitch();
    this.snackbar.dismiss();
  }

  preview() {
    this.previewTx = true;
    this.unlockAndSend();
  }

  send() {
    this.previewTx = false;
    this.unlockAndSend();
  }

  unlockAndSend() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.previewButton.resetState();
    this.sendButton.resetState();

    if (this.form.get('wallet').value.encrypted) {
      this.dialog.open(PasswordDialogComponent).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      this.createTransaction();
    }
  }

  addDestination() {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.push(this.createDestinationFormGroup());
  }

  removeDestination(index) {
    const destinations = this.form.get('destinations') as FormArray;
    destinations.removeAt(index);
  }

  setShareValue(event) {
    this.autoShareValue = parseFloat(event.value).toFixed(2);
  }

  toggleOptions(event) {
    event.stopPropagation();
    event.preventDefault();

    this.autoOptions = !this.autoOptions;
  }

  setAutoHours(event) {
    this.autoHours = event.checked;
    this.form.get('destinations').updateValueAndValidity();

    if (!this.autoHours) {
      this.autoOptions = false;
    }
  }

  private fillForm() {
    this.addresses = this.formData.form.wallet.addresses;

    ['wallet', 'addresses', 'changeAddress'].forEach(name => {
      this.form.get(name).setValue(this.formData.form[name]);
    });

    for (let i = 0; i < this.formData.form.destinations.length - 1; i++) {
      this.addDestination();
    }

    this.destControls.forEach((destControl, i) => {
      ['address', 'coins', 'hours'].forEach(name => {
        destControl.get(name).setValue(this.formData.form.destinations[i][name]);
      });
    });

    if (this.formData.form.hoursSelection.type === 'auto') {
      this.autoShareValue = this.formData.form.hoursSelection.share_factor;
      this.autoHours = true;
    } else {
      this.autoHours = false;
    }

    this.autoOptions = this.formData.form.autoOptions;
  }

  addressCompare(a, b) {
    return a.address === b.address;
  }

  get destControls() {
    return (this.form.get('destinations') as FormArray).controls;
  }

  private validateDestinations() {
    if (!this.form || !Array.isArray(this.form.get('addresses').value)) {
      return { Required: true };
    }

    const invalidInput = this.destControls.find(control => {
      const checkControls = ['coins'];

      if (!this.autoHours) {
        checkControls.push('hours');
      }

      return checkControls.map(name => {
        const value = control.get(name).value !== undefined
          ? control.get(name).value.replace(' ', '=')
          : '';

        if (isNaN(value) || value.trim() === '') {
          return true;
        }

        if (parseFloat(value) <= 0) {
          return true;
        }

        if (name === 'coins') {
          const parts = value.split('.');

          if (parts.length === 2 && parts[1].length > 6) {
            return true;
          }
        } else if (name === 'hours') {
          if (value < 1 || parseInt(value, 10) !== parseFloat(value)) {
            return true;
          }
        }

        return false;
      }).find(e => e === true);
    });

    if (invalidInput) {
      return { Invalid: true };
    }

    const coins = this.form.get('addresses').value.reduce((a, b) => a + b.coins, 0);
    const hours = this.form.get('addresses').value.reduce((a, b) => a + b.hours, 0);
    const destinationsCoins = this.destControls.reduce((a, b) => a + parseFloat(b.value.coins), 0);
    const destinationsHours = this.destControls.reduce((a, b) => a + parseInt(b.value.hours, 10), 0);

    if (destinationsCoins > coins || destinationsHours > hours) {
      return { Invalid: true };
    }

    return null;
  }

  private createDestinationFormGroup() {
    return this.formBuilder.group({
      address: '',
      coins: '',
      hours: '',
    });
  }

  private createTransaction(passwordDialog?: any) {
    if (passwordDialog) {
      passwordDialog.close();
    }

    if (this.previewTx) {
      this.previewButton.setLoading();
      this.sendButton.setDisabled();
    } else {
      this.sendButton.setLoading();
      this.previewButton.setDisabled();
    }

    this.walletService.createTransaction(
      this.form.get('wallet').value,
      this.form.get('addresses').value.map(addr => addr.address),
      this.destinations,
      this.hoursSelection,
      this.form.get('changeAddress').value ? this.form.get('changeAddress').value : null,
      passwordDialog ? passwordDialog.password : null,
    )
      .toPromise()
      .then(transaction => {
        if (!this.previewTx) {
          return this.walletService.injectTransaction(transaction.encoded).toPromise();
        }

        this.onFormSubmitted.emit({
          form: {
            wallet: this.form.get('wallet').value,
            addresses: this.form.get('addresses').value,
            changeAddress: this.form.get('changeAddress').value,
            destinations: this.destinations,
            hoursSelection: this.hoursSelection,
            autoOptions: this.autoOptions,
          },
          amount: this.destinations.reduce((a, b) => a + parseFloat(b.coins), 0),
          to: this.destinations.map(d => d.address),
          transaction,
        });
      })
      .then(() => {
        this.sendButton.setSuccess();
        this.resetForm();

        setTimeout(() => {
          this.sendButton.resetState();
        }, 3000);
      })
      .catch(error => {
        showSnackbarError(this.snackbar, error);

        this.previewButton.resetState().setEnabled();
        this.sendButton.resetState().setEnabled();
      });
  }

  private resetForm() {
    this.form.get('wallet').setValue('', { emitEvent: false });
    this.form.get('addresses').setValue([]);
    this.form.get('changeAddress').setValue('');

    while (this.destControls.length > 0) {
      (this.form.get('destinations') as FormArray).removeAt(0);
    }

    this.addDestination();

    this.autoHours = true;
    this.autoOptions = false;
    this.autoShareValue = '0.5';
  }

  private get destinations() {
    return this.destControls.map(destControl => {
      const destination = {
        address: destControl.get('address').value,
        coins: destControl.get('coins').value,
      };

      if (!this.autoHours) {
        destination['hours'] = destControl.get('hours').value;
      }

      return destination;
    });
  }

  private get hoursSelection() {
    let hoursSelection = {
      type: 'manual',
    };

    if (this.autoHours) {
      hoursSelection = <any> {
        type: 'auto',
        mode: 'share',
        share_factor: this.autoShareValue,
      };
    }

    return hoursSelection;
  }
}
