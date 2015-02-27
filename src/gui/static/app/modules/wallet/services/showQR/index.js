'use strict';

require('angular-qrcode');

module.exports = angular.module('skycoin.wallet.services.showQR', [
  'monospaced.qrcode'
])
.factory('showQR', require('./service'))
.controller('ModalCtrl', require('./controller'));

module.name = 'skycoin.wallet.services.showQR';
