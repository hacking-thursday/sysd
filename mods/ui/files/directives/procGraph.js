app

.controller("ProcGraphCtrl"
	, [       "$scope", "Loading", "sysd"
	, function($scope,   Loading,   sysd) {

	$scope.Loading = Loading;
	$scope.procData = {};
	$scope.curProc = null;
	$scope.curSock = null;
	$scope.curJson = null;
	$scope.jsonData = "";

	$scope.selectProc = function(proc) {
		$scope.curJson = $scope.curProc = proc;
	};

	$scope.selectSock = function(sock) {
		$scope.curJson = $scope.curSock = sock;
	};

	$scope.$watch("curJson", function() {
		$scope.jsonData = JSON.stringify($scope.curJson, undefined, 4);
	});

	sysd.getProcGraph().then(function(data) {
		$scope.procData = data;
		for (var pid in data.process) {
			// init with first process data
			$scope.selectProc(data.process[pid]);
			break;
		}
	});

}])

;
