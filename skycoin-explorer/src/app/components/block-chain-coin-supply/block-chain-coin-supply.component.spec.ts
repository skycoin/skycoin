import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BlockChainCoinSupplyComponent } from './block-chain-coin-supply.component';

describe('BlockChainCoinSupplyComponent', () => {
  let component: BlockChainCoinSupplyComponent;
  let fixture: ComponentFixture<BlockChainCoinSupplyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BlockChainCoinSupplyComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BlockChainCoinSupplyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
