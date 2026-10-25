import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Subject, of, throwError } from 'rxjs';

import { HardwareWalletComponent, HwPageStates, readFeatures } from './hardware-wallet.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('HardwareWalletComponent', () => {
  let component: HardwareWalletComponent;
  let fixture: ComponentFixture<HardwareWalletComponent>;
  let hw: any;
  let msgBar: jasmine.SpyObj<MsgBarService>;
  let dialog: any;

  const features = {
    label: 'mine', fw_major: 1, fw_minor: 8, fw_patch: 0,
    initialized: true, pin_protection: true, needs_backup: false, bootloader_mode: false,
  };

  function create(connected = true, feat: any = features) {
    hw = {
      walletConnectedAsyncEvent: new Subject<boolean>(),
      getDeviceConnected: jasmine.createSpy('getDeviceConnected').and.returnValue(of(connected)),
      getFeatures: jasmine.createSpy('getFeatures').and.returnValue(of({ rawResponse: feat })),
      wipe: jasmine.createSpy('wipe').and.returnValue(of({})),
      changeLabel: jasmine.createSpy('changeLabel').and.returnValue(of({})),
      recoverMnemonic: jasmine.createSpy('recoverMnemonic').and.returnValue(of({})),
      changePin: jasmine.createSpy('changePin').and.returnValue(throwError(() => ({ translatableErrorMsg: 'hardware-wallet.errors.incorrect-pin' }))),
    };
    msgBar = jasmine.createSpyObj('MsgBarService', ['hide', 'showDone', 'showError']);
    dialog = { open: jasmine.createSpy('open') };
    TestBed.configureTestingModule({
      declarations: [HardwareWalletComponent, MockTranslatePipe],
      imports: [FormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        { provide: HwWalletService, useValue: hw },
        { provide: MsgBarService, useValue: msgBar },
        { provide: CustomMatDialogService, useValue: dialog },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(HardwareWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  it('reads the device features', () => {
    expect(readFeatures(features)).toEqual({
      label: 'mine', firmware: '1.8.0', initialized: true, pinProtection: true, needsBackup: false, bootloaderMode: false,
    });
  });

  it('shows a connected device', () => {
    create();
    expect(component.state).toBe(HwPageStates.Ready);
    expect(component.device!.label).toBe('mine');
    // Change PIN, remove PIN, the two seed checks and wipe.
    expect(fixture.nativeElement.querySelectorAll('.-actions app-button').length).toBe(5);
    expect(fixture.nativeElement.querySelector('#hw-label')).toBeTruthy();
  });

  it('says when no device is connected', () => {
    create(false);
    expect(component.state).toBe(HwPageStates.Disconnected);
    expect(hw.getFeatures).not.toHaveBeenCalled();
  });

  it('offers to generate or restore a seed on a device without one', () => {
    create(true, { ...features, initialized: false });
    // Generate and restore, with 12 and 24 words, and nothing else.
    expect(fixture.nativeElement.querySelectorAll('.-actions app-button').length).toBe(4);
    expect(fixture.nativeElement.querySelector('#hw-label')).toBeNull();
  });

  it('checks the seed without changing the device', () => {
    create();
    component.checkSeed(24);
    expect(hw.recoverMnemonic).toHaveBeenCalledWith(24, true);
    expect(msgBar.showDone).toHaveBeenCalledWith('hardware-wallet.device.seed-matches');
  });

  it('shows the device error and refreshes', () => {
    create();
    hw.getFeatures.calls.reset();
    component.changePin();
    expect(hw.changePin).toHaveBeenCalledWith(true);
    expect(msgBar.showError).toHaveBeenCalledWith('hardware-wallet.errors.incorrect-pin');
    expect(hw.getFeatures).toHaveBeenCalled();
  });

  it('changes the label only when it differs', () => {
    create();
    component.changeLabel();
    expect(hw.changeLabel).not.toHaveBeenCalled();
    component.label = ' new ';
    component.changeLabel();
    expect(hw.changeLabel).toHaveBeenCalledWith('new');
  });

  it('wipes only after confirmation', () => {
    create();
    dialog.open.and.returnValue({ afterClosed: () => of(false) });
    component.wipe();
    expect(hw.wipe).not.toHaveBeenCalled();
    dialog.open.and.returnValue({ afterClosed: () => of(true) });
    component.wipe();
    expect(hw.wipe).toHaveBeenCalled();
  });
});
