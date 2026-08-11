import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { UntypedFormBuilder } from '@angular/forms';

import { ScanAddressesComponent } from './scan-addresses.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { MockTranslatePipe, MockWalletService } from '../../../../utils/test-mocks';

describe('ScanAddressesComponent', () => {
  let component: ScanAddressesComponent;
  let fixture: ComponentFixture<ScanAddressesComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ ScanAddressesComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        UntypedFormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: {} },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ScanAddressesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
