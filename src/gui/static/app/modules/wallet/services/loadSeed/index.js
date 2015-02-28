'use strict';

var mod = 'skycoin.wallet.services.loadSeed';

module.exports = angular.module(mod, [])
.factory(mod + '.loadSeedService', require('./service'))
.controller(mod + 'ModalCtrl', require('./controller.js'));

module.name = mod;
