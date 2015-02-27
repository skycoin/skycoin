'use strict';

var OpenQR = function ($modal, $log) {
  return function(wallet){
    var modalInstance = $modal.open({
      template: require('./qrModal.jade'),
      controller: 'QRInstanceCtrl',
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

OpenQR.$inject = ['$modal', '$log'];
module.exports = OpenQR;
