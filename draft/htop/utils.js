// string repeat function
// 'x'.repeat( 100 )
String.prototype.repeat = function( num )
{
    return new Array( num + 1 ).join( this );
}


function calculate_cpu_usage(idle_1, idle_2, total_1, total_2) {
    return 100 * ((total_2 - total_1) - (idle_2 - idle_1)) / (total_2 - total_1);
}
