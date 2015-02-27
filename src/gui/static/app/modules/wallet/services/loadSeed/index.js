'use strict';

module.exports = angular.module('skycoin.wallet.services.loadSeed', [])
.factory('loadSeed', require('./service'))
.controller('ModalCtrl', require('./controller.js'));

module.name = 'skycoin.wallet.services.loadSeed';
