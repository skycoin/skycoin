'use strict';

var mod = 'skycoin.wallet.services';
var factory = function(showQR, loadSeed, update, create){
  return {
    showQR: showQR,
    loadSeed: loadSeed,
    update: update,
    create: create
  };
};

factory.$inject = [
  mod + '.showQR.showQRService',
  mod + '.loadSeed.loadSeedService',
  mod + '.update.updateService',
  mod + '.create.createService'
];

module.exports = angular.module(mod, [
  require('./showQR').name,
  require('./loadSeed').name,
  require('./update').name,
  require('./create').name
])
.factory('$wallet', factory);

module.name = mod;
