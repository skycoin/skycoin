'use strict';


// Declare app level module which depends on filters, and services
angular.module('skycoin', [
  'ngRoute',
  'ui.bootstrap',
  //'skycoin.filters',
  //'skycoin.services',
  //'skycoin.directives',
  'skycoin.controllers',
  'monospaced.qrcode'
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view1', {templateUrl: 'partials/partial1.html', controller: 'mainCtrl'});
  $routeProvider.when('/view2', {templateUrl: 'partials/partial2.html', controller: 'mainCtrl'});
  $routeProvider.otherwise({redirectTo: '/'});
}]).config(['$locationProvider', function($locationProvider){
    $locationProvider.html5Mode(true).hashPrefix('');
}]);


