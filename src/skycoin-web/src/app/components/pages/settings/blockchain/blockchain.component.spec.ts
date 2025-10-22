import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { BlockchainComponent } from './blockchain.component';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { BlockchainService } from '../../../../services/blockchain.service';
import { CoinService } from '../../../../services/coin.service';
import { MockTranslatePipe, MockCoinService, MockDateTimePipe, MockBlockchainService } from '../../../../utils/test-mocks';

describe('BlockchainComponent', () => {
  let component: BlockchainComponent;
  let fixture: ComponentFixture<BlockchainComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        BlockchainComponent,
        MockDateTimePipe,
        MockTranslatePipe
      ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: CoinService, useClass: MockCoinService }
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BlockchainComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
