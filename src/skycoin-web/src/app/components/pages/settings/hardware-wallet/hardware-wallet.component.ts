import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subscription, mergeMap, of } from 'rxjs';

import { HwWalletService, OperationResult } from '../../../../services/hw-wallet.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { showConfirmationModal } from '../../../../utils';
import { OperationError } from '../../../../utils/operation-error';
import { MessageIcons } from '../../../layout/hardware-wallet/hw-message/hw-message.component';

export enum HwPageStates {
  Checking = 'checking',
  Disconnected = 'disconnected',
  Error = 'error',
  Ready = 'ready',
  Working = 'working',
}

/** What the device reports about itself, from /features. */
export interface HwDeviceInfo {
  label: string;
  firmware: string;
  initialized: boolean;
  pinProtection: boolean;
  needsBackup: boolean;
  bootloaderMode: boolean;
}

/**
 * Manages the connected Skywallet: its seed, PIN, label and backup. The device asks
 * for confirmation of each operation on its own screen.
 */
@Component({
    selector: 'app-hardware-wallet',
    templateUrl: './hardware-wallet.component.html',
    styleUrls: ['./hardware-wallet.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class HardwareWalletComponent implements OnInit, OnDestroy {
  states = HwPageStates;
  icons = MessageIcons;
  state = HwPageStates.Checking;
  device: HwDeviceInfo | null = null;
  error = '';
  label = '';
  readonly maxLabelLength = HwWalletService.maxLabelLength;

  private operation?: Subscription;
  private connection?: Subscription;

  constructor(
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
    private dialog: CustomMatDialogService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.connection = this.hwWalletService.walletConnectedAsyncEvent.subscribe(() => {
      if (this.state !== HwPageStates.Working) {
        this.refresh();
      }
      this.changeDetectorRef.markForCheck();
    });
    this.refresh();
  }

  ngOnDestroy() {
    this.connection?.unsubscribe();
    this.operation?.unsubscribe();
    this.msgBarService.hide();
  }

  refresh() {
    this.state = HwPageStates.Checking;
    this.operation?.unsubscribe();
    this.operation = this.hwWalletService.getDeviceConnected().pipe(
      mergeMap(connected => connected ? this.hwWalletService.getFeatures() : of(null)),
    ).subscribe(result => {
      if (!result) {
        this.device = null;
        this.state = HwPageStates.Disconnected;
      } else {
        this.device = readFeatures(result.rawResponse);
        this.label = this.device.label;
        this.state = HwPageStates.Ready;
      }
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.device = null;
      this.error = errorText(err);
      this.state = HwPageStates.Error;
      this.changeDetectorRef.markForCheck();
    });
  }

  generateSeed(wordCount: number) {
    this.run(this.hwWalletService.generateMnemonic(wordCount), 'hardware-wallet.device.seed-generated');
  }

  restoreSeed(wordCount: number) {
    this.run(this.hwWalletService.recoverMnemonic(wordCount, false), 'hardware-wallet.device.seed-restored');
  }

  checkSeed(wordCount: number) {
    this.run(this.hwWalletService.recoverMnemonic(wordCount, true), 'hardware-wallet.device.seed-matches');
  }

  backup() {
    this.run(this.hwWalletService.backup(), 'hardware-wallet.device.backed-up');
  }

  changePin() {
    this.run(this.hwWalletService.changePin(this.device!.pinProtection), 'hardware-wallet.device.pin-changed');
  }

  removePin() {
    this.run(this.hwWalletService.removePin(), 'hardware-wallet.device.pin-removed');
  }

  changeLabel() {
    const label = this.label.trim();
    if (!label || label === this.device!.label) {
      return;
    }
    this.run(this.hwWalletService.changeLabel(label), 'hardware-wallet.device.label-changed');
  }

  wipe() {
    showConfirmationModal(this.dialog, {
      headerText: 'hardware-wallet.device.wipe-title',
      text: 'hardware-wallet.device.wipe-warning',
      checkboxText: 'hardware-wallet.device.wipe-check',
      confirmButtonText: 'hardware-wallet.device.wipe-button',
      cancelButtonText: 'confirmation.cancel-button',
      redTitle: true,
    }).afterClosed().subscribe(confirmed => {
      if (confirmed) {
        this.run(this.hwWalletService.wipe(), 'hardware-wallet.device.wiped');
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  private run(operation: Observable<OperationResult>, doneText: string) {
    this.msgBarService.hide();
    this.state = HwPageStates.Working;
    this.operation?.unsubscribe();
    this.operation = operation.subscribe(() => {
      this.msgBarService.showDone(doneText);
      this.refresh();
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.msgBarService.showError(errorText(err));
      this.refresh();
      this.changeDetectorRef.markForCheck();
    });
    this.changeDetectorRef.markForCheck();
  }
}

export function readFeatures(features: any): HwDeviceInfo {
  return {
    label: features.label || '',
    firmware: [features.fw_major, features.fw_minor, features.fw_patch].join('.'),
    initialized: !!features.initialized,
    pinProtection: !!features.pin_protection,
    needsBackup: !!features.needs_backup,
    bootloaderMode: !!features.bootloader_mode,
  };
}

function errorText(err: any): string {
  const e = err as OperationError;

  return (e && (e.translatableErrorMsg || e.originalServerErrorMsg)) || (err && err.message) || 'hardware-wallet.errors.generic-error';
}
