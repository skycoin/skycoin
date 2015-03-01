'use strict';

// @ngInject
var openQR = function ($modal, $log) {
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

module.exports = openQR;
