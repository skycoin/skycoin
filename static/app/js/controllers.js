'use strict';

/* Controllers */

angular.module('myApp.controllers', []).
  controller('MyCtrl1', [function() {

  }])
  .controller('MyCtrl2', [function() {

  }])

.controller('TabsDemoCtrl', ['$scope',
  function($scope) {

	$scope.tabs = [
	    { title:"Transactions", content:"Dynamic content 2", disabled: false },
      { title:"Addresses", content:"Dynamic content 2", disabled: false }
	  ];

  $scope.alertMe = function() {
    setTimeout(function() {
      //alert("You've selected the alert tab!");
    });
  };

  $scope.navType = 'pills';
  
}]);