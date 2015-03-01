'use strict';

require('./controllers');

// Declare app level module which depends on filters, and services
module.exports = angular.module('skycoin.home', [
  'ui.router',
  'ui.bootstrap',
  'skycoin.controllers',
  'skycoin.wallet.services'
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
