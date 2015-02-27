'use strict';

module.exports = angular.module('skycoin.wallet.services', [
  require('./showQR').name,
  require('./loadSeed').name
])
.factory('$wallet', function(showQR, loadSeed){
  return {
    showQR: showQR,
    loadSeed: loadSeed
  };
});

module.name = 'skycoin.wallet';
