app

.factory("sysdProcGraph"
	, [       "$q", "sysd"
	, function($q,   sysd) {

	sysd.supportProcGraph = false;

	var depApis = ["process/resource", "network/socket", "network/ifce", "sysfs"];

	sysd.checkProcGraph = function() {
		var res = {
			apis: []
		};

		// check all necessry api supported
		angular.forEach(depApis, function(api) {
			if (!sysd.api.get[api]) {
				res.apis.push(api);
				res.errMsg = "some dependent apis unsupported in this platform";
			}
		});
		sysd.supportProcGraph = res.supported = !res.apis[0];
		return res;
	};

	sysd.getProcGraph = function() {
		var deferred = $q.defer();
		var check = sysd.checkProcGraph();

		if (!check.supported) {
			deferred.reject(check);
			return deferred.promise;
		}

		var deferall = [];
		angular.forEach(depApis, function(api) {
			deferall.push(sysd.api.get[api]());
		});
		$q.all(deferall).then(function(res) {
			process_data = res[0].processes;
			socket_data = res[1];
                        ifce_data = res[2];
                        hw_data = res[3];

			// 編列 process 的資料
			result_process = {}
			for ( var key in process_data ){
				pid = key;

				row = {};
				row["pid"] = pid;
				row["cmdline"] = process_data[key]["cmdline"];
				row["status"] = process_data[key]["status"];
				row["socket"] = [];

				fd_ary = process_data[key]["fd"];
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

                        result_process_tree = {}
                        {
                            function check_and_init( p_id, pids2, result_process ){
                                if( p_id == "0" ){ return; }

                                if ( !pids2[p_id] ){
                                    pids2[p_id] = { "obj": result_process[p_id] };

                                    pids2[p_id]["obj"]["pid"] = p_id;
                                    pids2[p_id]["obj"]["children"] = [];
                                    //pids2[p_id]["obj"]["parent"] = null;
                                }
                            }
                            pids = Object.keys( result_process );
                            pids2 = {};
                            for( pid = pids.pop() ; pid ; pid = pids.pop() ){
                                ppid = result_process[pid]["status"]["ppid"];

                                check_and_init( ppid, pids2, result_process );
                                check_and_init(  pid, pids2, result_process );

                                pids2[pid]["obj"] = result_process[pid];
                                if( ppid != "0" ){
                                    pids2[ppid]["obj"]["children"].push( pids2[pid]["obj"] );
                                }
                            }
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

                        // 編列 interface 的資料
                        result_ifce = {}
			for ( i=0; i<ifce_data.length; i++ ){
                            ifce_name = ifce_data[i]["Name"];
                            result_ifce[ifce_name] = ifce_data[i];
                        }

                        // 編列 hardware 的資料
                        result_hw = {}
                        for ( var ifce_name in hw_data["class"]["net"] ){
                                var row = {};
                                row["device_path"] = hw_data["class"]["net"][ifce_name][1];
                                result_hw[ifce_name] = row;
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

			deferred.resolve({
				process: result_process,
				socket: result_socket,
				ifce: result_ifce
			});
		}, function(res) {
			deferred.reject(res);
		});
		return deferred.promise;
	};

	return sysd;
}])

;
