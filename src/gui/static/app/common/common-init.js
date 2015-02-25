'use strict';

function commonInit($rootScope, $state) {
	$rootScope.$state = $state;
	
	// Set bodyClasses, pageTitle, and pageDescription on state change (ui-router)
    $rootScope.$on('$stateChangeSuccess', function(event, toState){
		if ( angular.isDefined( toState.data.pageTitle ) ) {
			$rootScope.pageTitle = toState.data.pageTitle;
			$rootScope.pageDescription = toState.data.pageDescription;
			$rootScope.bodyClasses = toState.data.moduleClasses + ' ' + toState.data.pageClasses;
		}
	});

	// Make sure the page scrolls to the top on all state transitions
	$rootScope.$on('$viewContentLoaded', function(){
		if (document.readyState === 'complete') {
			window.scrollTo(0, 0);
		}
	});

    // Proper Regex Pattern for email input form validation
	$rootScope.emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

}

commonInit.$inject = ['$rootScope', '$state'];
module.exports = commonInit;