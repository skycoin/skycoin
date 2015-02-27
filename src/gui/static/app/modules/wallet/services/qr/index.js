'use strict';

module.exports = angular.module('skycoin.wallet.services.qr', [])
.factory('OpenQR', require('./qrService'))
.controller('QRInstanceCtrl', require('./qrInstanceCtrl'));

module.name = 'skycoin.wallet.services.qr';
