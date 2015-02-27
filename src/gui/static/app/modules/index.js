'use strict';

module.exports = angular.module('modules', [
		require('./home').name,
		require('./qr').name
  ]
)
.controller('MainCtrl', require('./MainController'));
