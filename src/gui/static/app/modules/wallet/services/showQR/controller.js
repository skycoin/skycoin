'use strict';

// Controller naming conventions should start with an uppercase letter
// @ngInject
function ModalCtrl($http, $scope, $modalInstance, wallet) {

  $scope.address = wallet.entries[0].address;

  $scope.ok = function () {
    $modalInstance.close();
  };
}
module.exports = ModalCtrl;
