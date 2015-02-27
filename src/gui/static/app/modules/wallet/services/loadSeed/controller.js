'use strict';

// Controller naming conventions should start with an uppercase letter
// @ngInject
function ModalCtrl($http, $scope, $modalInstance) {

  $scope.wallet = {};
  $scope.wallet.name = '';
  $scope.wallet.new = '';

  $scope.ok = function () {

    console.log('New wallet called');
    var xsrf = {
      name:$scope.wallet.name,
      seed:$scope.wallet.seed
    };
    console.log('xsrf: ', xsrf);
    $http({
      method: 'POST',
      url: '/wallet/create',
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
      console.log('Load wallet response: ');
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
