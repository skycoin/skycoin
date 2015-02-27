'use strict';

var loadSeed = function ($modal, $log) {
  return function(wallet){
    var modalInstance = $modal.open({
      template: require('./modal.jade'),
      controller: require('./controller'),
      resolve: {
        wallet: function () {
          return wallet;
        }
      }
    });

    modalInstance.result.then(function () {
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };
};

loadSeed.$inject = ['$modal', '$log'];
module.exports = loadSeed;
