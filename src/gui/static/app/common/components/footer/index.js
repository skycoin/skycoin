'use strict';

module.exports = angular.module('common.components.commonFooter', [])
.directive('commonFooter', function () {
  return {
    template: require('./common-footer.html'),
    restrict: 'EA',
    replace: true
  };
});
