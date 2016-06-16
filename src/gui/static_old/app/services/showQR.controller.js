(function() {
  angular.module('app.services')
  .controller('ShowQRController', Controller);

  // @ngInject
  function Controller($http, $scope, $modalInstance, wallet) {

    $scope.address = wallet.entries[0].address;

    $scope.ok = function () {
      $modalInstance.close();
    };
  }

})();
