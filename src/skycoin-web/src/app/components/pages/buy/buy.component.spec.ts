import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';

import { BuyComponent } from './buy.component';
import { PurchaseService } from '../../../services/purchase.service';
import { MockTranslatePipe, MockPurchaseService, MockTellerStatusPipe, MockDateTimePipe, MockCustomMatDialogService } from '../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';

describe('BuyComponent', () => {
  let component: BuyComponent;
  let fixture: ComponentFixture<BuyComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [
        BuyComponent,
        MockTellerStatusPipe,
        MockDateTimePipe,
        MockTranslatePipe
      ],
      imports: [
        MatCardModule,
        MatDividerModule,
        MatIconModule,
        MatListModule
      ],
      providers: [
        { provide: PurchaseService, useClass: MockPurchaseService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BuyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
