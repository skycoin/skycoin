'use strict';

// Controller naming conventions should start with an uppercase letter
// @ngInject
function ModalCtrl($http, $scope, $modalInstance, wallet) {

  $scope.wallet = {};
  $scope.wallet.name = wallet.name;
  $scope.wallet.filename = wallet.filename;

  $scope.ok = function () {
    console.log('New wallet called');
    var xsrf = {
      name:$scope.wallet.name,
      id:$scope.wallet.filename
    };
    console.log('xsrf: ', xsrf);
    $http({
      method: 'POST',
      url: '/wallet/update',
      headers: {'Content-Type': 'application/x-www-form-urlencoded'},
      transformRequest: function(obj) {
        var str = [];
        for(var p in obj){
          str.push(encodeURIComponent(p) + '=' + encodeURIComponent(obj[p]));
        }
        return str.join('&');
      },
      data: xsrf
    }).success(function(response){
      console.log('Update wallet response: ');
      console.dir(response);
      //$scope.loadWallets();
    });

    $modalInstance.close();
  };

  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };
}
module.exports = ModalCtrl;
