'use strict';

module.exports = angular.module('modules',
	[
		require('./home').name
	])
	.controller('MainCtrl', require('./MainController'));