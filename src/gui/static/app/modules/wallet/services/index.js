'use strict';

var mod = 'skycoin.wallet.services';
var factory = function(showQR, loadSeed, update){
  return {
    showQR: showQR,
    loadSeed: loadSeed,
    update: update
  };
};

factory.$inject = [
  mod + '.showQR.showQRService',
  mod + '.loadSeed.loadSeedService',
  mod + '.update.updateService'
];

module.exports = angular.module(mod, [
  require('./showQR').name,
  require('./loadSeed').name,
  require('./update').name
])
.factory('$wallet', factory);

module.name = mod;
