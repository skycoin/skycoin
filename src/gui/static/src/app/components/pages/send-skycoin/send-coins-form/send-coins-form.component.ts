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
import { SelectAddressComponent } from '../../../layout/select-address/select-address';
import { BigNumber } from 'bignumber.js';
import { ConfirmationData } from '../../../../app.datatypes';
import { BlockchainService } from '../../../../services/blockchain.service';
import { showConfirmationModal } from '../../../../utils';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { ChangeNoteComponent } from '../send-preview/transaction-info/change-note/change-note.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { MultipleDestinationsDialogComponent } from '../../../layout/multiple-destinations-dialog/multiple-destinations-dialog.component';
import { FormSourceSelectionComponent, AvailableBalanceData, SelectedSources, SourceSelectionModes } from '../form-parts/form-source-selection/form-source-selection.component';
import { FormDestinationComponent, Destination } from '../form-parts/form-destination/form-destination.component';
import { CopyRawTxComponent, CopyRawTxData } from '../offline-dialogs/implementations/copy-raw-tx.component';
import { DoubleButtonActive } from '../../../../components/layout/double-button/double-button.component';

@Component({
  selector: 'app-send-coins-form',
  templateUrl: './send-coins-form.component.html',
  styleUrls: ['./send-coins-form.component.scss'],
})
export class SendCoinsFormComponent implements OnInit, OnDestroy {
  public static lastShowForManualUnsignedValue = false;

  @ViewChild('formSourceSelection', { static: false }) formSourceSelection: FormSourceSelectionComponent;
  @ViewChild('formMultipleDestinations', { static: false }) formMultipleDestinations: FormDestinationComponent;
  @ViewChild('previewButton', { static: false }) previewButton: ButtonComponent;
  @ViewChild('sendButton', { static: false }) sendButton: ButtonComponent;
  @Input() formData: any;
  @Input() showSimpleForm: boolean;
  @Output() onFormSubmitted = new EventEmitter<any>();

  sourceSelectionModes = SourceSelectionModes;
  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  form: FormGroup;
  availableBalance = new AvailableBalanceData();
  selectedSources: SelectedSources;
  autoHours = true;
  autoOptions = false;
  autoShareValue = '0.5';
  previewTx: boolean;
  busy = false;
  showForManualUnsigned = SendCoinsFormComponent.lastShowForManualUnsignedValue;
  doubleButtonActive = DoubleButtonActive;

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
    this.msgBarService.hide();

    SendCoinsFormComponent.lastShowForManualUnsignedValue = this.showForManualUnsigned;
  }

  sourceSelectionChanged() {
    this.selectedSources = this.formSourceSelection.selectedSources;
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

  changeFormType(value: DoubleButtonActive) {
    if ((value === DoubleButtonActive.LeftButton && !this.showForManualUnsigned) || (value === DoubleButtonActive.RightButton && this.showForManualUnsigned)) {
      return;
    }

    if (value === DoubleButtonActive.RightButton) {
      const confirmationData: ConfirmationData = {
        text: 'send.unsigned-confirmation',
        headerText: 'confirmation.header-text',
        confirmButtonText: 'confirmation.confirm-button',
        cancelButtonText: 'confirmation.cancel-button',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.showForManualUnsigned = true;
        }
      });
    } else {
      this.showForManualUnsigned = false;
    }
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

    if (!this.showForManualUnsigned && this.selectedSources.wallet.encrypted && !this.selectedSources.wallet.isHardware && !this.previewTx) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.selectedSources.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.createTransaction(passwordDialog);
        });
    } else {
      if (this.previewTx || this.showForManualUnsigned || !this.selectedSources.wallet.isHardware) {
        this.createTransaction();
      } else {
        this.showBusy();
        this.processingSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.selectedSources.wallet.addresses[0].address).subscribe(
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
            };
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

    let selectedAddresses: string[];

    if (!this.showForManualUnsigned) {
      selectedAddresses = this.selectedSources.addresses && this.selectedSources.addresses.length > 0 ?
        this.selectedSources.addresses.map(addr => addr.address) : null;
    } else {
      selectedAddresses = this.selectedSources.manualAddresses;
    }

    const selectedOutputs = this.selectedSources.unspentOutputs && this.selectedSources.unspentOutputs.length > 0 ?
      this.selectedSources.unspentOutputs.map(addr => addr.hash) : null;

    const destinations = this.formMultipleDestinations.getDestinations(!this.autoHours);

    this.processingSubscription = this.walletService.createTransaction(
      this.selectedSources.wallet,
      selectedAddresses ? selectedAddresses : this.selectedSources.wallet.addresses.map(address => address.address),
      selectedOutputs,
      destinations,
      this.hoursSelection,
      this.form.get('changeAddress').value ? this.form.get('changeAddress').value : null,
      passwordDialog ? passwordDialog.password : null,
      this.previewTx || !this.selectedSources.wallet,
    ).subscribe(transaction => {
      if (passwordDialog) {
        passwordDialog.close();
      }

      const note = this.form.value.note.trim();
      if (!this.previewTx) {
        if (!this.showForManualUnsigned) {
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
          const data: CopyRawTxData = {
            rawTx: transaction.encoded,
            isUnsigned: true,
          };

          const config = new MatDialogConfig();
          config.width = '566px';
          config.data = data;

          this.dialog.open(CopyRawTxComponent, config).afterClosed().subscribe(() => {
            this.resetState();

            const confirmationData: ConfirmationData = {
              text: 'offline-transactions.copy-tx.reset-confirmation',
              headerText: 'confirmation.header-text',
              confirmButtonText: 'confirmation.confirm-button',
              cancelButtonText: 'confirmation.cancel-button',
            };

            showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
              if (confirmationResult) {
                this.resetForm();
                this.msgBarService.showDone('offline-transactions.copy-tx.reset-done', 4000);
              }
            });
          });
        }
      } else {
        let amount = new BigNumber('0');
        destinations.map(destination => amount = amount.plus(destination.coins));
        this.onFormSubmitted.emit({
          form: {
            wallet: this.selectedSources.wallet,
            addresses: this.selectedSources.addresses,
            manualAddresses: this.selectedSources.manualAddresses,
            changeAddress: this.form.get('changeAddress').value,
            destinations: destinations,
            hoursSelection: this.hoursSelection,
            autoOptions: this.autoOptions,
            allUnspentOutputs: this.formSourceSelection.unspentOutputsList,
            outputs: this.selectedSources.unspentOutputs,
            currency: this.formMultipleDestinations.currentlySelectedCurrency,
            note: note,
          },
          amount: amount,
          to: destinations.map(d => d.address),
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

  private resetState() {
    this.busy = false;
    this.navbarService.enableSwitch();
    this.previewButton.resetState().setEnabled();
    this.sendButton.resetState().setEnabled();
  }
}
