import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { of } from 'rxjs';
import { ActivatedRoute } from '@angular/router';

import { OutputsComponent } from './outputs.component';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { CoinService } from '../../../../services/coin.service';
import { MockTranslatePipe, MockSpendingService, MockCoinService, MockCustomMatDialogService } from '../../../../utils/test-mocks';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

describe('OutputsComponent', () => {
  let component: OutputsComponent;
  let fixture: ComponentFixture<OutputsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OutputsComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: ActivatedRoute, useValue: { queryParams: of({}) } },
        { provide: SpendingService, useClass: MockSpendingService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OutputsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
