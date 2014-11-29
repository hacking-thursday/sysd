app

.controller("IndexCtrl"
	, [       "$scope", "$mdBottomSheet", "$mdToast", "$translate", "localStorageService", "sysd", "sysdProcGraph", "Loading"
	, function($scope,   $mdBottomSheet,   $mdToast,   $translate,   localStorageService,   sysd,   sysdProcGraph,   Loading) {

	window.sysd = sysd;
	sysd.host = localStorageService.get("sysdhost") || "127.0.0.1";
	sysd.port = +localStorageService.get("sysdport") || 8080;
	$scope.sysd = sysd;
	$scope.Loading = Loading;
	$scope.config = {
		host: sysd.host,
		port: sysd.port
	};
	$scope.apiReq = {
		method: "",
		path: ""
	};
	$scope.apiResData = "";

	$scope.showLangBottomSheet = function($event) {
		// show config view
		$scope.apiResData = "";

		$mdBottomSheet.show({
			templateUrl: "directives/langSelect.html",
			controller: "langSelectCtrl",
			targetEvent: $event
		}).then(function(clickedItem) {
			$translate.use(clickedItem.lang);
		});
	};

	$scope.api = function(path) {
		Loading.add("api");
		$scope.apiReq.method = "GET";
		$scope.apiReq.path = "/" + path;
		sysd.api.get[path]().then(function(data) {
			$scope.apiResData = JSON.stringify(data, undefined, 4);
		}, function() {
			$mdToast.show({
				template: "<md-toast>" +
					"<div translate>Execute API failed</div>" +
					"<div style='margin-left: 1em;'>" + path + "</div>" +
					"</md-toast>",
				hideDelay: 5000,
				position: "top left"
			});
		}).finally(function() {
			Loading.del("api");
		});
	};

	$scope.regapis = function() {
		Loading.add("regapis");
		sysd.regapis().then(function() {
			// success
		}, function(res) {
			if (sysd.host == "127.0.0.1" && sysd.port == 8080) {
				return;
			}
			$mdToast.show({
				template: "<md-toast>" +
					"<div translate>Connect to sysd server failed</div>" +
					"<div translate style='margin-left: 1em;'>Please check configuration</div>" +
					"<div style='margin-left: 1em;'>" + sysd.host + ":" + sysd.port + "</div>" +
					"</md-toast>",
				hideDelay: 5000,
				position: "top left"
			});
		}).finally(function() {
			Loading.del("regapis");
		});
	};

	$scope.saveConfig = function() {
		sysd.host = $scope.config.host;
		sysd.port = $scope.config.port;
		localStorageService.set("sysdhost", sysd.host);
		localStorageService.set("sysdport", sysd.port);
		$scope.regapis();
	};

	$scope.resetConfig = function() {
		sysd.host = $scope.config.host = "127.0.0.1";
		sysd.port = $scope.config.port = 8080;
		localStorageService.clearAll()
		$scope.regapis();
	};

	$scope.regapis();

}])

;
