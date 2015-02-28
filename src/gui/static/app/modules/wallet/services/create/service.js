'use strict';

// @ngInject
var create = function ($modal, $log) {
  return function(loadWallets){
    var modalInstance = $modal.open({
      template: require('./modal.jade'),
      controller: require('./controller'),
      resolve: {
        loadWallets: function () {
          return loadWallets;
        }
      }
    });

    modalInstance.result.then(function () {
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };
};

module.exports = create;
