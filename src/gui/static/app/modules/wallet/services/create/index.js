'use strict';

var mod = 'skycoin.wallet.services.create';
module.exports = angular.module(mod, [])
.factory(mod + '.createService', require('./service'))
.controller(mod + '.ModalCtrl', require('./controller'));

module.name = mod;
