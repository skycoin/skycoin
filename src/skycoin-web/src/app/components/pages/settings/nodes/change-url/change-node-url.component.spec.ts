import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { FormBuilder } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { TranslateService } from '@ngx-translate/core';

import { ChangeNodeURLComponent } from './change-node-url.component';
import { MockTranslatePipe, MockCoinService, MockMsgBarService } from '../../../../../utils/test-mocks';
import { CoinService } from '../../../../../services/coin.service';
import { MsgBarService } from '../../../../../services/msg-bar.service';

describe('ChangeNodeURLComponent', () => {
  let component: ChangeNodeURLComponent;
  let fixture: ComponentFixture<ChangeNodeURLComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChangeNodeURLComponent, MockTranslatePipe ],
      imports: [ HttpClientModule ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        FormBuilder,
        { provide: CoinService, useClass: MockCoinService },
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: {} },
        {
          provide: TranslateService,
          useValue: jasmine.createSpyObj('TranslateService', ['instant'])
        },
        { provide: MsgBarService, useClass: MockMsgBarService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChangeNodeURLComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
