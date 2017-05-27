"use strict";
var testing_1 = require('@angular/core/testing');
var block_chain_coin_supply_component_1 = require('./block-chain-coin-supply.component');
describe('BlockChainCoinSupplyComponent', function () {
    var component;
    var fixture;
    beforeEach(testing_1.async(function () {
        testing_1.TestBed.configureTestingModule({
            declarations: [block_chain_coin_supply_component_1.BlockChainCoinSupplyComponent]
        })
            .compileComponents();
    }));
    beforeEach(function () {
        fixture = testing_1.TestBed.createComponent(block_chain_coin_supply_component_1.BlockChainCoinSupplyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });
    it('should create', function () {
        expect(component).toBeTruthy();
    });
});
//# sourceMappingURL=block-chain-coin-supply.component.spec.js.map