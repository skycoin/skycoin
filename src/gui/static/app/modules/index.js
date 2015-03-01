'use strict';

require('./home');

module.exports = angular.module('modules', [
  'skycoin.home',
  require('./wallet').name
]);
