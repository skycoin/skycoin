'use strict';

// Controller naming conventions should start with an uppercase letter
function ModalCtrl($http, $scope, $modalInstance, wallet) {

  $scope.address = wallet.entries[0].address;

  $scope.ok = function () {
    $modalInstance.close();
  };
}
// $inject is necessary for minification. See http://bit.ly/1lNICde for explanation.
ModalCtrl.$inject = ['$http', '$scope', '$modalInstance', 'wallet'];

module.exports = ModalCtrl;
