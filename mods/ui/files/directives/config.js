app

.controller("ConfigCtrl"
	, [       "$scope", "localStorageService", "sysd"
	, function($scope,   localStorageService,   sysd) {

	$scope.saveConfig = function() {
		sysd.host = $scope.config.host;
		sysd.port = $scope.config.port;
		localStorageService.set("sysdhost", sysd.host);
		localStorageService.set("sysdport", sysd.port);
		$scope.regapis();
	};

	$scope.resetConfig = function() {
		sysd.host = $scope.config.host = "127.0.0.1";
		sysd.port = $scope.config.port = 8;
		localStorageService.clearAll()
		$scope.regapis();
	};

}])

;
