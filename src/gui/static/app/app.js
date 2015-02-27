'use strict';

require('angular');

module.exports = angular.module('skycoin',
	[
		require('./common/common.js').name,
		require('./modules').name
	])
	.config(require('./appConfig'))
	.constant('version', require('../package.json').version)
	.run(require('./common/common-init.js'));
