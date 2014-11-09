app

.factory("sysd"
	, [       "$http", "$q"
	, function($http,   $q) {

	function sysd() {
		var $scope = this;
		$scope.host = "127.0.0.1";
		$scope.port = 8080;
		$scope.api = {
			get: {}
		};
		$scope.apis = [];
		return $scope;
	}

	sysd.prototype.makeApiUrl = function(path) {
		var $scope = this;
		return "http://" + $scope.host + ":" + $scope.port + path;
	};

	sysd.prototype.callapi = function(method, path) {
		var $scope = this;
		var deferred = $q.defer();
		var url = $scope.makeApiUrl(path);
		$http({
			url: url,
			method: method
		})
		.success(function(data, status, headers, config) {
			deferred.resolve(data, status, headers, config);
		})
		.error(function(data, status, headers, config) {
			deferred.reject({
				data: data,
				status: status,
				headers: headers,
				config: config
			});
		});
		return deferred.promise;
	};

	sysd.prototype.regapis = function() {
		var $scope = this;
		var deferred = $q.defer();
		$scope.api = {
			get: {}
		};
		$scope.apis = [];
		$scope.callapi("GET", "/apilist")
			.then(function(apis) {
				angular.forEach(apis, function(path) {
					path = path.replace(/^\//, "");
					if (path.substr(0,3) == "ui/") {
						return
					}
					$scope.apis.push(path);

					$scope.api.get[path] = function() {
						return $scope.callapi("GET", "/" + path);
					};
				});
				$scope.apis.sort();
				deferred.resolve(apis);
			}, function(res) {
				deferred.reject(res);
			});
		return deferred.promise;
	};

	return new sysd();
}])

;
