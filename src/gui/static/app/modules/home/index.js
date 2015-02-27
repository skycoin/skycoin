'use strict';

require('./controllers');
require('angular-qrcode');

// Declare app level module which depends on filters, and services
module.exports = angular.module('skycoin', [
  'ui.router',
  'ui.bootstrap',
  'skycoin.controllers',
  'monospaced.qrcode'
])
.config([
  '$stateProvider', '$urlRouterProvider',
  function($stateProvider, $urlRouterProvider) {
    $stateProvider.state('view1', {
      url: '/view1',
      template: require('./partial1.html'),
      controller: 'mainCtrl'
    });
    $stateProvider.state('view2', {
      url: '/view2',
      template: require('./partial2.html'),
      controller: 'mainCtrl'
    });
    $urlRouterProvider.otherwise('/');
  }
])
.config(['$locationProvider', function($locationProvider){
    $locationProvider.html5Mode(true).hashPrefix('');
  }
]);
