'use strict';

// Controller naming conventions should start with an uppercase letter
function QRInstanceCtrl($http, $scope, $modalInstance, wallet) {

  $scope.address = wallet.entries[0].address;
  $scope.qro = {};
  $scope.qro.fm = wallet.entries[0].address;

  $scope.$watch('qro.label', function() {
    $scope.qro.new = 'skycoin:' + $scope.address.address;// + '?' + 'label=' + $scope.qro.label; //+ '&message=' + $scope.qro.message;
  });

  $scope.$watch('qro.message', function() {
    $scope.qro.new = 'skycoin:' + $scope.address.address;// + '?' + 'label=' + $scope.qro.label; //+ '&message=' + $scope.qro.message;
  });


  $scope.ok = function () {
    $modalInstance.close($scope.selected.item);
  };

  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };

}
// $inject is necessary for minification. See http://bit.ly/1lNICde for explanation.
QRInstanceCtrl.$inject = ['$http', '$scope', '$modalInstance', 'wallet'];

module.exports = QRInstanceCtrl;
