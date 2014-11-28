app

.factory("sysdProcGraph"
	, [       "$q", "sysd"
	, function($q,   sysd) {

	sysd.getProcGraph = function() {
		var deferred = $q.defer();
		var apis = ["process/resource", "network/socket"];
		var err = {
			apis: []
		};

		// check all necessry api supported
		angular.forEach(apis, function(api) {
			if (!sysd.api.get[api]) {
				err.apis.push(api);
			}
		});
		if (err.apis[0]) {
			deferred.reject({
				errMsg: "some dependent apis unsupported in this platform",
				apis: err.apis
			});
			return deferred.promise;
		}

		var deferall = [];
		angular.forEach(apis, function(api) {
			deferall.push(sysd.api.get[api]());
		});
		$q.all(deferall).then(function(res) {
			// TODO: parse res
			deferred.resolve(res);
		}, function(res) {
			deferred.reject(res);
		});
		return deferred.promise;
	};

	return sysd;
}])

;
