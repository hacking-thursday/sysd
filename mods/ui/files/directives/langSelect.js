app

.controller("langSelectCtrl"
	, [       "$scope", "$mdBottomSheet"
	, function($scope,   $mdBottomSheet) {

	$scope.items = [
		{ name: "English", lang: "en_us", icon: 'hangout' },
		{ name: "正體中文", lang: "zh_tw", icon: 'hangout' }
	];

	$scope.listItemClick = function($index) {
		var clickedItem = $scope.items[$index];
		$mdBottomSheet.hide(clickedItem);
	};

}])

;
