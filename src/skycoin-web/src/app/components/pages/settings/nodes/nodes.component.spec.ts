import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { NodesComponent } from './nodes.component';
import { MockTranslatePipe, MockCoinService, MockCustomMatDialogService } from '../../../../utils/test-mocks';
import { CoinService } from '../../../../services/coin.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';

describe('NodesComponent', () => {
  let component: NodesComponent;
  let fixture: ComponentFixture<NodesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NodesComponent, MockTranslatePipe ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: CoinService, useClass: MockCoinService },
        { provide: CustomMatDialogService, useClass: MockCustomMatDialogService }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
