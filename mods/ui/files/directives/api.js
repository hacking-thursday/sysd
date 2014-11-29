app

.controller("ApiCtrl"
	, [       "$scope", "$routeParams", "$mdToast", "localStorageService", "Loading", "sysd"
	, function($scope,   $routeParams,   $mdToast,   localStorageService,   Loading,   sysd) {

	$scope.Loading = Loading;
	$scope.apiname = $routeParams.apiname;

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

	if ($scope.apiname) {
		$scope.api($scope.apiname);
	}

}])

;
