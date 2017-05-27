"use strict";
var testing_1 = require('@angular/core/testing');
var coin_supply_service_1 = require('./coin-supply.service');
describe('CoinSupplyService', function () {
    beforeEach(function () {
        testing_1.TestBed.configureTestingModule({
            providers: [coin_supply_service_1.CoinSupplyService]
        });
    });
    it('should ...', testing_1.inject([coin_supply_service_1.CoinSupplyService], function (service) {
        expect(service).toBeTruthy();
    }));
});
//# sourceMappingURL=coin-supply.service.spec.js.map