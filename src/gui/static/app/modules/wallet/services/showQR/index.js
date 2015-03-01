'use strict';

require('angular-qrcode');

var mod = 'skycoin.wallet.services.showQR';

module.exports = angular.module(mod, [
  'monospaced.qrcode'
])
.factory(mod + '.showQRService', require('./service'))
.controller(mod + '.ModalCtrl', require('./controller'));

module.name = mod;
