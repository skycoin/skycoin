'use strict';

/* Controllers */

angular.module('skycoin.controllers', [])

.controller('mainCtrl', ['$scope','$http',
  function($scope,$http) {

  	$scope.loadWallet = function(wallet){
	  var data = {WalletName:wallet};
	  console.log('wallet is loading:' + wallet);
      $http.post('/api/loadWallet', JSON.stringify(data)).success(function(response){
        console.dir(response);
        $scope.loadedWallet = response;
      });
	 }

	 console.log('local storage wallet is ' + localStorage.loadedWallet)
	 $scope.loadWallet(localStorage.loadedWallet);

	 $scope.saveWallet = function(){
	  var data = {Address:$scope.newAddress};
      $http.post('/api/saveWallet', JSON.stringify(data)).success(function(response){
        console.dir(response);
        $scope.loadedWalletName = response;
        localStorage.loadedWallet = response.replace(/"/g, "");
      });
	 }

	 $scope.newAddress = function(){
	  	$http.get('/api/newAddress').success(function(response) {
	      console.dir(response);
	      $scope.newAddress = response.replace(/"/g, "");
	      $scope.saveWallet();
	    });
	 }






}]);