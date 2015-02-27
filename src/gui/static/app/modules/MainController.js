'use strict';
// Controller naming conventions should start with an uppercase letter
// @ngInject
function MainCtrl($rootScope, $scope) {
	$scope.test = null;
	console.log('Up and running!');
}

module.exports = MainCtrl;
