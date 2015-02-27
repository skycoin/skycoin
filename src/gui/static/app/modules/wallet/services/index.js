'use strict';

module.exports = angular.module('skycoin.wallet.services', [
  require('./qr').name
])
.factory('$wallet', function(OpenQR){
  return {
    qr: OpenQR
  };
});

module.name = 'skycoin.wallet';
