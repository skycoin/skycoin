'use strict';
// Controller naming conventions should start with an uppercase letter
function HomeViewCtrl($scope) {
	$scope.testVar = 'We are up and running using a required module!';
}

// $inject is necessary for minification. See http://bit.ly/1lNICde for explanation.
HomeViewCtrl.$inject = ['$scope'];
module.exports = HomeViewCtrl;