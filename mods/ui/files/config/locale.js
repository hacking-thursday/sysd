app

.config(["$translateProvider", function($translateProvider) {

	var zh_tw = {
		"Save": "保存",
		"Reset": "清除",
		"After configure sysd host and port, click left menu to get response from Sysd server": "設定完 Sysd 的連線方式後，點擊左側選單，從 Sysd 伺服器取得資料",
		"Connect to sysd server failed": "連線到 Sysd 伺服器時發生問題",
		"Please check configuration": "請檢查相關設定",
		"Execute API failed": "執行 API 發生問題"
	};

	$translateProvider
	.translations("zh_tw", zh_tw)
	.translations("zh_TW", zh_tw);

	var defTrans = {};
	angular.forEach($translateProvider.translations("zh_tw"), function(v, key) {
		if (key.match(/{{.*}}/)) {
			defTrans[key] = key;
		}
	});

	$translateProvider
	.translations("en", defTrans)
	.registerAvailableLanguageKeys(["en", "zh_tw"], {
		"en*": "en",
		"zh*": "zh_tw"
	})
	.useStorage("localStorageService")
	.determinePreferredLanguage();

}])

;
