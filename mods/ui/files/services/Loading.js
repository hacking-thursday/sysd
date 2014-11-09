app

.factory("Loading", function() {
	function Loading() {
		var $scope = this;
		$scope.val = {};
		return $scope;
	}

	Loading.prototype.add = function(name) {
		var $scope = this;
		if ($scope.val[name]) {
			$scope.val[name]++;
		} else {
			$scope.val[name] = 1;
		}
		return $scope;
	};

	Loading.prototype.del = function(name) {
		var $scope = this;
		if ($scope.val[name]) {
			$scope.val[name]--;
		} else {
			$scope.val[name] = 0;
		}
		return $scope;
	};

	Loading.prototype.isLoading = function(name) {
		var $scope = this;
		return !!$scope.val[name];
	};

	return new Loading();
})

;
