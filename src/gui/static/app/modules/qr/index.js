'use strict';

module.exports = angular.module('skycoin.qr', [])
.factory('OpenQR', require('./qrService'))
.controller('QRInstanceCtrl', require('./qrInstanceCtrl'));
