import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { MatMenuModule, MatIconModule, MatTooltipModule } from '@angular/material';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { TopBarComponent } from './top-bar.component';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { CoinService } from '../../../../services/coin.service';
import { MockTranslatePipe, MockBalanceService, MockCoinService, MockLanguageService, MockCustomMatDialogService } from '../../../../utils/test-mocks';
import { LanguageService } from '../../../../services/language.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

describe('TopBarComponent', () => {
  let component: TopBarComponent;
  let fixture: ComponentFixture<TopBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TopBarComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      imports: [
        MatMenuModule,
        MatIconModule,
        MatTooltipModule,
        RouterTestingModule
      ],
      providers: [
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: LanguageService, useClass: MockLanguageService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TopBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
