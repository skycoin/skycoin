import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { UntypedFormBuilder } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { MatSelectModule } from '@angular/material/select';
import { NO_ERRORS_SCHEMA, Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';

import { AddDepositAddressComponent } from './add-deposit-address.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { PurchaseService } from '../../../../services/purchase.service';
import { MockTranslatePipe, MockPurchaseService, MockWalletService } from '../../../../utils/test-mocks';

describe('AddDepositAddressComponent', () => {
  let component: AddDepositAddressComponent;
  let fixture: ComponentFixture<AddDepositAddressComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AddDepositAddressComponent, MockTranslatePipe ],
      imports: [ MatSelectModule ],
      providers: [
        UntypedFormBuilder,
        { provide: WalletService, useClass: MockWalletService },
        { provide: PurchaseService, useClass: MockPurchaseService },
        { provide: MatDialogRef, useValue: {} }
      ],
      schemas: [ NO_ERRORS_SCHEMA ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddDepositAddressComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
