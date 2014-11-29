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
			window.data = [];
			process_data = res[0].processes;
			socket_data = res[1];

			// 編列 process 的資料
			result_process = {}
			for ( var key in process_data ){
				pid = key;

				row = {};
				row["pid"] = pid;
				row["socket"] = [];

				fd_ary = process_data[key];
				for( var fd in fd_ary ){
					fd_key = fd_ary[fd][0];
					fd_val = fd_ary[fd][1];
					matches = fd_val.match( /socket:\[(\d+)\]/ );
					if ( matches ) {
						var itm = {};
						itm["type"] = "socket";
						itm["inode"] = matches[1];

						row["socket"].push( itm );
					}

				}


				result_process[pid] = row;
			}

			// 編列 socket 的資料
			result_socket = {}
			for ( i=0; i<socket_data.length; i++ ){
				var itm = socket_data[i];
				if ( itm.socket_type && ( itm.inode || itm.Inode ) ){
					var inode = itm.inode;
					if ( itm.Inode ){
						inode = itm.Inode;
					}
					result_socket[inode] = itm;
				}
			}

			// 重新 scan 一次 process 的列表，並將 socket 的資料用物件連結取代
			for( var key in result_process ){
				var itm = result_process[key];
				for( var j=0; j < itm["socket"].length; j++ ){
					var inode = itm["socket"][j]["inode"];
					if ( result_socket[inode] ){
						result_process[key]["socket"][j] = result_socket[inode];
					}
				}
			}

			// 注意，目前先暫時將資料存在 window.data 裡下，方便 debug
			window.data = {};
			window.data.process = result_process;
			window.data.socket = result_socket;

			deferred.resolve(res);
		}, function(res) {
			deferred.reject(res);
		});
		return deferred.promise;
	};

	return sysd;
}])

;
