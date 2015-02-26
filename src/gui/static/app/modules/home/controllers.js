'use strict';

/* Controllers */

module.exports = angular.module('skycoin.controllers', [])

.controller('mainCtrl', ['$scope','$http', '$modal', '$log', '$timeout',
  function($scope,$http,$modal,$log,$timeout) {
    $scope.addresses = [];

    $scope.tab = {};

    $scope.sendTab = function(){
      $scope.tab.sendActive = true;
    };

    $scope.getProgress = function(){
      $http.get('/blockchain/progress').success(function(response){
        console.log('Block chain progress: ');
        console.dir(response);
        $scope.progress = (parseInt(response.current,10)+1) / parseInt(response.Highest,10) * 100;

      });
    };

    $scope.refreshBalances = function() {
      $scope.loadWallets();
      $timeout($scope.refreshBalances, 15000);
    };

    $scope.getProgress();
    $timeout($scope.refreshBalances, 15000);

    $scope.loadWallets = function(){
      $http.post('/wallets').success(function(response){
        console.log('Loading wallets');
        console.dir(response);
        //$scope.loadedWallet = response;
        $scope.wallets = response;
        /*for(var i=0;i<response.length;i++){
          if(!$scope.addresses[i]) $scope.addresses[i] = {};
          $scope.addresses[i].address = response[i].address;
        }*/
        for(var i=0;i<response.length;i++){
          $scope.checkBalance(i,response[i].address);
        }
      });
    };

    console.log('local storage wallet is ' + localStorage.loadedWallet);

    $scope.loadWallets();

    $scope.saveWallet = function(){
      var data = {Addresses:$scope.addresses};
      $http.post('/wallets/save', JSON.stringify(data)).success(function(response){
        console.log('Wallet Save: ');
        console.dir(response);
        $scope.loadedWalletName = response;
        localStorage.loadedWallet = response.replace(/"/g, '');
      });
    };

    $scope.newWallet = function(){
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
        $scope.loadWallets();
      });
    };


    $scope.sendDisable = true;
    $scope.readyDisable = false;

    $scope.ready = function(){
      $scope.readyDisable = !$scope.readyDisable;
      $scope.sendDisable = !$scope.sendDisable;
    };

    $scope.clearSend = function(){
      $scope.sendDisable = true;
      $scope.readyDisable = false;
      $scope.spend.amount = '';
      $scope.spend.address = '';
    };

    $scope.pendingTable = [];

    //reset local storage
    localStorage.setItem('historyTable',JSON.stringify([]));

    $scope.historyTable = JSON.parse(localStorage.getItem('historyTable'));
    console.log('localStorage.history');
    console.dir(JSON.parse(localStorage.getItem('historyTable')));

    $scope.spend = function(spend){
      $scope.sendDisable = true;
      $scope.readyDisable = true;
      $timeout($scope.clearSend, 1000);
      $scope.pendingTable.push(spend);
      var xsrf = {
        id:spend.id,
        coins:spend.amount*1000000,
        fee:1,
        hours:1,
        address:spend.address
      };
      console.log('spend xsrf is ' , xsrf);
      $scope.historyTable.push({address:spend.address,amount:spend.amount});
      localStorage.setItem('historyTable',JSON.stringify($scope.historyTable));
      console.dir($scope.historyTable);
      $http({
        method: 'POST',
        url: '/wallet/spend',
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
        console.log('wallet spend is ');
        console.dir(response);
        $scope.loadWallets();
      });
    };

    $scope.checkBalance = function(wI, address){
      console.log('Checking bal of: ');
      console.dir(wI);
      console.log('Checking address: ');
      console.dir(address);
      var xsrf = {addr:address};
      $http({
        method: 'POST',
        url: '/wallet/balance',
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
        console.log('Check balance: ');
        console.dir(response);
        $scope.wallets[wI].balance = response.confirmed.coins / 1000000;
      });
    };

    $scope.mainBackUp = function(){
    };

    $scope.openQR = function (wallet) {
      var modalInstance = $modal.open({
        template: require('./qr-modal.html'),
        controller: 'qrInstanceCtrl',
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

    $scope.openLoadWallet = function (wallet) {

      var modalInstance = $modal.open({
        template: require('./loadWalletModal.html'),
        controller: 'loadWalletInstanceCtrl',
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

    $scope.updateWallet = function (wallet) {

      var modalInstance = $modal.open({
        template: require('./updateWalletModal.html'),
        controller: 'updateWalletInstanceCtrl',
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

  }
])


.controller('qrInstanceCtrl', ['$http', '$scope', '$modalInstance', 'wallet',
  function($http, $scope, $modalInstance, wallet) {

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
}])

.controller('loadWalletInstanceCtrl', ['$http', '$scope', '$modalInstance', 'wallet',
  function($http, $scope, $modalInstance) {

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
}])

.controller('updateWalletInstanceCtrl', ['$http', '$scope', '$modalInstance', 'wallet',
  function($http, $scope, $modalInstance, wallet) {

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
}]);
