'use strict';

module.exports = angular.module('skycoin.wallet.services.showQR', [])
.factory('showQR', require('./service'))
.controller('ModalCtrl', require('./controller'));

module.name = 'skycoin.wallet.services.showQR';
