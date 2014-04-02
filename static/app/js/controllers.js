'use strict';

/* Controllers */

angular.module('skycoin.controllers', [])

.controller('mainCtrl', ['$scope','$http', '$modal', '$log', '$timeout',
  function($scope,$http,$modal,$log,$timeout) {
  	$scope.addresses = [];

  	$scope.tab = {};

  	$scope.sendTab = function(){
      $scope.tab.sendActive = true;
	 };


  	$scope.getProgress = function(){
      $http.get('/blockchain/progress').success(function(response){
      	console.log('Block chain progress: ')
      	console.dir(response)
        $scope.progress = (parseInt(response.current,10)+1) / parseInt(response.Highest,10) * 100;

      });
	 }

	$scope.refreshBalances = function() {
	    $scope.loadWallets();
	    $timeout($scope.refreshBalances, 15000);
    }

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
	 }

	 console.log('local storage wallet is ' + localStorage.loadedWallet)

	 $scope.loadWallets();

	 $scope.saveWallet = function(){
	  var data = {Addresses:$scope.addresses};
      $http.post('/wallets/save', JSON.stringify(data)).success(function(response){
      	console.log('Wallet Save: ');
        console.dir(response);
        $scope.loadedWalletName = response;
        localStorage.loadedWallet = response.replace(/"/g, "");
      });
	 }

	 $scope.newAddress = function(){
	  	$http.get('/wallet').success(function(response) {
	  		conole.log('New address: ')
		     console.dir(response);
		     //$scope.addresses.push(response.address);
		     $scope.loadWallets();
		     //$scope.addresses.push(response.replace(/"/g, ""));
		     //$scope.saveWallet();
	    });
	 }

	 $scope.sendDisable = true;
	 $scope.readyDisable = false;

	 $scope.ready = function(){
	 	$scope.readyDisable = !$scope.readyDisable;
	 	$scope.sendDisable = !$scope.sendDisable;
	 }

	 $scope.clearSend = function(){
	 	$scope.sendDisable = true;
	 	$scope.readyDisable = false;
	 	$scope.spend.amount = '';
	 	$scope.spend.address = '';
	 }

	 $scope.pendingTable = [];


	 //reset local storage
	 localStorage.setItem('historyTable',JSON.stringify([]));


	 $scope.historyTable = JSON.parse(localStorage.getItem('historyTable'));
	 console.log('localStorage.history')
	 console.dir(JSON.parse(localStorage.getItem('historyTable')))

	 $scope.spend = function(addr){
	 	$scope.sendDisable = true;
	 	$scope.readyDisable = true;
	 	$timeout($scope.clearSend, 1000);
	 	$scope.pendingTable.push(addr);
	 	var xsrf = {dst:addr.address,
	 				coins:addr.amount*1000000,
	 				fee:1,
	 				hours:1}
	 	$scope.historyTable.push({address:addr.address,amount:addr.amount});
	 	localStorage.setItem('historyTable',JSON.stringify($scope.historyTable));
	 	console.dir($scope.historyTable);
		$http({
		    method: 'POST',
		    url: '/wallet/spend',
		    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
		    transformRequest: function(obj) {
		        var str = [];
		        for(var p in obj)
		        str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
		        return str.join("&");
		    },
		    data: xsrf
			}).success(function(response){
		  	 	console.log('wallet spend is ')
		        console.dir(response);
		        $scope.loadWallets();
	      });
	 }

	 $scope.checkBalance = function(wI, address){
	 	console.log('Checking bal of: ');
	 	console.dir(wI);
	 	console.log('Checking address: ');
	 	console.dir(address)
	 	var xsrf = {addr:address}
		$http({
		    method: 'POST',
		    url: '/wallet/balance',
		    headers: {'Content-Type': 'application/x-www-form-urlencoded'},
		    transformRequest: function(obj) {
		        var str = [];
		        for(var p in obj)
		        str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
		        return str.join("&");
		    },
		    data: xsrf
			}).success(function(response){
				console.log("Check balance: ");
				console.dir(response);
		        $scope.wallets[wI].balance = response.confirmed.coins / 1000000;
	      });
	 }

	 $scope.mainBackUp = function(){

	 }

	 $scope.openQR = function (address) {

      var modalInstance = $modal.open({
        templateUrl: 'qrModalContent.html',
        controller: "qrInstanceCtrl",
        resolve: {
          address: function () {
            return address;
          }
        }
      });

      modalInstance.result.then(function () {
      }, function () {
        $log.info('Modal dismissed at: ' + new Date());
      });
    };

}])


.controller('qrInstanceCtrl', ['$http', '$scope', '$modalInstance', 'address',
  function($http, $scope, $modalInstance, address) {

  $scope.address = address;
  $scope.qro = {};
  $scope.qro.fm = address;

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
}]);