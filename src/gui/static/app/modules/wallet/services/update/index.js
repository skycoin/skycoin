'use strict';

var mod = 'skycoin.wallet.services.update';
module.exports = angular.module(mod, [])
.factory(mod + '.updateService', require('./service'))
.controller(mod + '.ModalCtrl', require('./controller'));

module.name = mod;
