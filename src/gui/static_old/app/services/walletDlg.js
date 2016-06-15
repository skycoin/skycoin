(function() {
  angular.module('app.services')

  .factory('walletDlg', function($modal, $log) {
    var create = function(loadWallets){
      var modalInstance = $modal.open({
        templateUrl: 'services/createWallet.html',
        controller: 'CreateWalletController',
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

    var update = function(wallet){
      var modalInstance = $modal.open({
        templateUrl: 'services/updateWallet.html',
        controller: 'UpdateWalletController',
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

    var showQR = function(wallet){
      var modalInstance = $modal.open({
        templateUrl: 'services/showQR.html',
        controller: 'ShowQRController',
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

    var loadSeed = function(wallet){
      var modalInstance = $modal.open({
        templateUrl: 'services/loadSeed.html',
        controller: 'LoadSeedController',
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

    return {
      create: create,
      update: update,
      showQR: showQR,
      loadSeed: loadSeed
    };
  });

})();
