function getCPUUsage($scope, $http, $interval) {
    var refresh_time = 2;  // second

    var cpu_total_first = new Array();
    var cpu_idle_first = new Array();
    var cpu_total_second = new Array();
    var cpu_idle_second = new Array();

    $interval(function(){
        // first read 
        $http.get('http://127.0.0.1:8/procstat').
            success(function(data, status, headers, config) {
                $.each(data['CPUs'], function( index, value ){
                    cpu_idle_first[index] = value['Idle'];
                    cpu_total_first[index] = value['User'] + value['Nice'] + value['System'] + value['Idle'] + value['Iowait'] + value['Irq'] + value['Softirq'];
                });
            }).
            error(function(data, status, headers, config) {
                // log error
            });

        // second read
        setTimeout(function() {
            $http.get('http://127.0.0.1:8/procstat').
                success(function(data, status, headers, config) {
                    $.each(data['CPUs'], function( index, value ){
                        cpu_idle_second[index] = value['Idle'];
                        cpu_total_second[index] = value['User'] + value['Nice'] + value['System'] + value['Idle'] + value['Iowait'] + value['Irq'] + value['Softirq'];
                    });

                    $scope.cpus = [];
                    for(i=0; i<cpu_total_first.length; i++) {
                        percent = cacular_cpu_usage(cpu_idle_first[i], cpu_idle_second[i], cpu_total_first[i], cpu_total_second[i]);
                        $scope.cpus.push({'id': i, 'percent': percent, 'bar': '|'.repeat(Math.round(percent / 100 * 140))});
                    }
                }).
                error(function(data, status, headers, config) {
                    // log error
                });
        }, 500);
    }, refresh_time * 1000);

    // $scope.cpus = [
    //     {
    //         id: 1,
    //         percent: 20
    //     },
    //     {
    //         id: 2,
    //         percent: 24
    //     },
    //     {
    //         id: 3,
    //         percent: 19
    //     },
    //     {
    //         id: 4,
    //         percent: 21
    //     }
    // ];
}
