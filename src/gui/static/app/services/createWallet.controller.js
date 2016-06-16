(function() {
  angular.module('app.services')
  .controller('CreateWalletController', Controller);

  // @ngInject
  function Controller($http, $scope, $modalInstance, loadWallets) {

    $scope.ok = function(){
      console.log('New wallet called');
      var xsrf = {name:''};
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
        console.log('New wallet response: ');
        console.dir(response);
        loadWallets();
      });
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  }

})();
