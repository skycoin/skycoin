import { SubscriptionLike } from 'rxjs';
import { first } from 'rxjs/operators';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormGroup, FormControl } from '@angular/forms';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { ButtonComponent } from '../../../layout/button/button.component';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { NavBarService } from '../../../../services/nav-bar.service';
import { SelectAddressComponent } from './select-address/select-address';
import { BigNumber } from 'bignumber.js';
import { ConfirmationData } from '../../../../app.datatypes';
import { BlockchainService } from '../../../../services/blockchain.service';
import { showConfirmationModal } from '../../../../utils';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { ChangeNoteComponent } from '../send-preview/transaction-info/change-note/change-note.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { MultipleDestinationsDialogComponent } from '../../../layout/multiple-destinations-dialog/multiple-destinations-dialog.component';
import { FormSourceSelectionComponent, AvailableBalanceData } from '../form-parts/form-source-selection/form-source-selection.component';
import { FormMultipleDestinationsComponent, Destination } from '../form-parts/form-multiple-destinations/form-multiple-destinations.component';

@Component({
  selector: 'app-send-form-advanced',
  templateUrl: './send-form-advanced.component.html',
  styleUrls: ['./send-form-advanced.component.scss'],
})
export class SendFormAdvancedComponent implements OnInit, OnDestroy {
  @ViewChild('formSourceSelection', { static: false }) formSourceSelection: FormSourceSelectionComponent;
  @ViewChild('formMultipleDestinations', { static: false }) formMultipleDestinations: FormMultipleDestinationsComponent;
  @ViewChild('previewButton', { static: false }) previewButton: ButtonComponent;
  @ViewChild('sendButton', { static: false }) sendButton: ButtonComponent;
  @Input() formData: any;
  @Output() onFormSubmitted = new EventEmitter<any>();

  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  form: FormGroup;
  availableBalance = new AvailableBalanceData();
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx: boolean;
  busy = false;

  private syncCheckSubscription: SubscriptionLike;
  private processingSubscription: SubscriptionLike;

  constructor(
    public blockchainService: BlockchainService,
    public walletService: WalletService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private navbarService: NavBarService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private changeDetector: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.navbarService.showSwitch('send.simple', 'send.advanced', DoubleButtonActive.RightButton);

    this.form = new FormGroup({}, this.validateForm.bind(this));
    this.form.addControl('changeAddress', new FormControl(''));
    this.form.addControl('note', new FormControl(''));

    if (this.formData) {
      setTimeout(() => this.fillForm());
    }
  }

  ngOnDestroy() {
    if (this.processingSubscription && !this.processingSubscription.closed) {
      this.processingSubscription.unsubscribe();
    }
    this.closeSyncCheckSubscription();
    this.navbarService.hideSwitch();
    this.msgBarService.hide();
  }

  sourceSelectionChanged() {
    this.availableBalance = this.formSourceSelection.availableBalance;
    if (this.formMultipleDestinations) {
      this.formMultipleDestinations.updateValuesAndValidity();
    }
    this.form.updateValueAndValidity();
  }

  multipleDestinationsChanged() {
    setTimeout(() => {
      this.form.updateValueAndValidity();
    });
  }

  preview() {
    this.previewTx = true;
    this.checkBeforeSending();
    this.changeDetector.detectChanges();
  }

  send() {
    this.previewTx = false;
    this.checkBeforeSending();
  }

  private checkBeforeSending() {
    if (!this.form.valid || this.previewButton.isLoading() || this.sendButton.isLoading()) {
      return;
    }

    this.closeSyncCheckSubscription();
    this.syncCheckSubscription = this.blockchainService.synchronized.pipe(first()).subscribe(synchronized => {
      if (synchronized) {
        this.prepareTransaction();
      } else {
        this.showSynchronizingWarning();
      }
    });
  }

  private showSynchronizingWarning() {
    const confirmationData: ConfirmationData = {
      text: 'send.synchronizing-warning',
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
    };

    showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.prepareTransaction();
      }
    });
  }

  private prepareTransaction() {
    this.msgBarService.hide();
    this.previewButton.resetState();
    this.sendButton.resetState();

    const selectedSources = this.formSourceSelection.selectedSources;

    if (selectedSources.wallet.encrypted && !selectedSources.wallet.isHardware && !this.previewTx) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: selectedSources.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (!selectedSources.wallet.isHardware || this.previewTx) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected(selectedSources.wallet.addresses[0].address).subscribe(
          () => this.createTransaction(),
          err => this.showError(getHardwareWalletErrorMsg(this.translate, err)),
        );
      }
    }
  }

  setShareValue(event) {
    this.autoShareValue = parseFloat(event.value).toFixed(2);
  }

  selectChangeAddress(event) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    this.dialog.open(SelectAddressComponent, config).afterClosed().subscribe(response => {
      if (response) {
        this.form.get('changeAddress').setValue(response);
      }
    });
  }

  openMultipleDestinationsPopup() {
    let currentString = '';

    const currentDestinations = this.formMultipleDestinations.getDestinations(!this.autoHours);
    currentDestinations.map(destControl => {
      if (destControl.address.trim().length > 0 ||
        destControl.coins.trim().length > 0 ||
        (!this.autoHours && destControl.hours.trim().length > 0)) {

          currentString += destControl.address.replace(',', '');
          currentString += ', ' + destControl.coins.replace(',', '');
          if (!this.autoHours) {
            currentString += ', ' + destControl.hours.replace(',', '');
          }
          currentString += '\r\n';
      }
    });

    if (currentString.length > 0) {
      currentString = currentString.substr(0, currentString.length - 1);
    }

    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = currentString;
    this.dialog.open(MultipleDestinationsDialogComponent, config).afterClosed().subscribe((response: string[][]) => {
      if (response) {
        if (response.length > 0) {
          this.autoHours = response[0].length === 2;

          const newDestinations: Destination[] = [];
          response.forEach((entry, i) => {
            const newDestination: Destination = {
              address: entry[0],
              coins: entry[1],
              originalAmount: null,
            }
            if (!this.autoHours) {
              newDestination.hours = entry[2];
            }

            newDestinations.push(newDestination);
          });

          this.formMultipleDestinations.setDestinations(newDestinations);
        } else {
          this.formMultipleDestinations.resetForm();
        }
      }
    });
  }

  toggleOptions(event) {
    event.stopPropagation();
    event.preventDefault();

    this.autoOptions = !this.autoOptions;
  }

  setAutoHours(event) {
    this.autoHours = event.checked;
    this.formMultipleDestinations.updateValuesAndValidity();

    if (!this.autoHours) {
      this.autoOptions = false;
    }
  }

  private fillForm() {
    this.formSourceSelection.fill(this.formData);
    this.formMultipleDestinations.fill(this.formData);

    ['changeAddress', 'note'].forEach(name => {
      this.form.get(name).setValue(this.formData.form[name]);
    });

    if (this.formData.form.hoursSelection.type === 'auto') {
      this.autoShareValue = this.formData.form.hoursSelection.share_factor;
      this.autoHours = true;
    } else {
      this.autoHours = false;
    }

    this.autoOptions = this.formData.form.autoOptions;
  }

  private validateForm() {
    if (!this.form) {
      return { Required: true };
    }

    if (!this.formSourceSelection || !this.formSourceSelection.valid || !this.formMultipleDestinations || !this.formMultipleDestinations.valid) {
      return { Invalid: true };
    }

    return null;
  }

  private showBusy() {
    if (this.previewTx) {
      this.previewButton.setLoading();
      this.sendButton.setDisabled();
    } else {
      this.sendButton.setLoading();
      this.previewButton.setDisabled();
    }
    this.busy = true;
    this.navbarService.disableSwitch();
  }

  private createTransaction(passwordDialog?: any) {
    this.showBusy();

    const selectedSources = this.formSourceSelection.selectedSources;

    const selectedAddresses = selectedSources.addresses && selectedSources.addresses.length > 0 ?
      selectedSources.addresses.map(addr => addr.address) : null;

    const selectedOutputs = selectedSources.unspentOutputs && selectedSources.unspentOutputs.length > 0 ?
      selectedSources.unspentOutputs.map(addr => addr.hash) : null;

      this.processingSubscription = this.walletService.createTransaction(
        selectedSources.wallet,
        selectedAddresses ? selectedAddresses : selectedSources.wallet.addresses.map(address => address.address),
        selectedOutputs,
        this.formMultipleDestinations.getDestinations(!this.autoHours),
        this.hoursSelection,
        this.form.get('changeAddress').value ? this.form.get('changeAddress').value : null,
        passwordDialog ? passwordDialog.password : null,
        this.previewTx,
      ).subscribe(transaction => {
        if (passwordDialog) {
          passwordDialog.close();
        }

        const note = this.form.value.note.trim();
        if (!this.previewTx) {
          this.processingSubscription = this.walletService.injectTransaction(transaction.encoded, note)
            .subscribe(noteSaved => {
              let showDone = true;
              if (note && !noteSaved) {
                this.msgBarService.showWarning(this.translate.instant('send.error-saving-note'));
                showDone = false;
              }

              this.showSuccess(showDone);
            }, error => this.showError(error));
        } else {
          let amount = new BigNumber('0');
          this.formMultipleDestinations.getDestinations(!this.autoHours).map(destination => amount = amount.plus(destination.coins));
          this.onFormSubmitted.emit({
            form: {
              wallet: selectedSources.wallet,
              addresses: selectedSources.addresses,
              changeAddress: this.form.get('changeAddress').value,
              destinations: this.formMultipleDestinations.getDestinations(!this.autoHours),
              hoursSelection: this.hoursSelection,
              autoOptions: this.autoOptions,
              allUnspentOutputs: this.formSourceSelection.unspentOutputsList,
              outputs: selectedSources.unspentOutputs,
              currency: this.formMultipleDestinations.currentlySelectedCurrency,
              note: note,
            },
            amount: amount,
            to: this.formMultipleDestinations.getDestinations(!this.autoHours).map(d => d.address),
            transaction,
          });
          this.busy = false;
          this.navbarService.enableSwitch();
        }
      }, error => {
        if (passwordDialog) {
          passwordDialog.error(error);
        }

        if (error && error.result) {
          this.showError(getHardwareWalletErrorMsg(this.translate, error));
        } else {
          this.showError(error);
        }
      });
  }

  private resetForm() {
    this.formSourceSelection.resetForm();
    this.formMultipleDestinations.resetForm();
    this.form.get('changeAddress').setValue('');
    this.form.get('note').setValue('');
    this.autoHours = true;
    this.autoOptions = false;
    this.autoShareValue = '0.5';
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

  private closeSyncCheckSubscription() {
    if (this.syncCheckSubscription) {
      this.syncCheckSubscription.unsubscribe();
    }
  }

  private showSuccess(showDone: boolean) {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.resetForm();

    if (showDone) {
      this.msgBarService.showDone('send.sent');
      this.sendButton.resetState();
    } else {
      this.sendButton.setSuccess();
      setTimeout(() => {
        this.sendButton.resetState();
      }, 3000);
    }
  }

  private showError(error) {
    this.busy = false;
    this.msgBarService.showError(error);
    this.navbarService.enableSwitch();
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }
}
