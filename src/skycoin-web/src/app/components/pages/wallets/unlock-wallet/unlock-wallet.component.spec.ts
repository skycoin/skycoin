import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { UntypedFormBuilder } from '@angular/forms';

import { UnlockWalletComponent } from './unlock-wallet.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { MockTranslatePipe, MockWalletService, MockMsgBarService } from '../../../../utils/test-mocks';
import { MsgBarService } from '../../../../services/msg-bar.service';

describe('UnlockWalletComponent', () => {
  let component: UnlockWalletComponent;
  let fixture: ComponentFixture<UnlockWalletComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ UnlockWalletComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        UntypedFormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: {} },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UnlockWalletComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
