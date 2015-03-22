'use strict';

// Services use camelCase for their names like Directives
// Factories have first letter capitalized like Controllers

module.exports = angular.module('common.services', [])
.factory('ServiceName', require('./serviceName.js'));
