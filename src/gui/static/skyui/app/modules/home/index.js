'use strict';
// Home View
module.exports = angular.module('home', [])
	.directive('homeView', require('./homeDirective'))
	.controller('HomeViewCtrl', require('./HomeController'));