'use strict';

module.exports = function homeDirective() {
	return {
		controller: 'HomeViewCtrl', // called from homeController.js
		template: require('./home.html'),
		restrict: 'EA',
		scope: true
	};
};