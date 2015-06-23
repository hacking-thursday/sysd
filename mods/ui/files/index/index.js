app

.config([     "$routeProvider"
	, function($routeProvider) {
		$routeProvider
			.when("/config", {
				templateUrl: "directives/config.html",
				controller: "ConfigCtrl"
			})
			.when("/api/:apiname*", {
				templateUrl: "directives/api.html",
				controller: "ApiCtrl"
			})
			.when("/procGraph", {
				templateUrl: "directives/procGraph.html",
				controller: "ProcGraphCtrl"
			})
			.otherwise({
				redirectTo: "/config"
			});
	}
])

.controller("IndexCtrl"
	, [       "$scope", "$mdBottomSheet", "$mdToast", "$translate", "sysd", "sysdProcGraph", "Loading"
	, function($scope,   $mdBottomSheet,   $mdToast,   $translate,   sysd,   sysdProcGraph,   Loading) {

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
		$mdBottomSheet.show({
			templateUrl: "directives/langSelect.html",
			controller: "langSelectCtrl",
			targetEvent: $event
		}).then(function(clickedItem) {
			$translate.use(clickedItem.lang);
		});
	};

	$scope.regapis = function() {
		Loading.add("regapis");
		sysd.regapis().then(function() {
			sysd.checkProcGraph();
		}, function(res) {
			if (sysd.host == "127.0.0.1" && sysd.port == 8) {
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

	$scope.regapis();

}])

;
