route - return the routing table information

## Description

| column          | description     |
| -------------   | :-------------  | 
| API Endpoint    | /route          |
| Method          | GET             |

## Parameters
| Field  | Type | Description                       |                  |
| -----  | ---- | -----------                       | -----            |
| pretty | int  | turn on pretty format mode or not | 0 = off , 1 = on |


## Return Values

Data Tree Structure:
```
.
`-- Root
    `-- DataRow
```

DataRow:

| Field       | Description | Sample Values |
| ----------  | ----------  | -----------   |
| Destination |             | 00000000      |
| Flags       |             | 0003          |
| Gateway     |             | 0114A8C0      |
| IRTT        |             | 0             |
| Iface       |             | wlan0         |
| MTU         |             | 0             |
| Mask        |             | 00000000      |
| Metric      |             | 303           |
| RefCnt      |             | 0             |
| Use         |             | 0             |
| Window      |             | 0             |

## Examples

### python
Example:
```
import urllib.request
import json

json_str = urllib.request.urlopen("http://127.0.0.1:8080/route").read(1000)
data_obj = json.loads(json_str.decode())
for row in data_obj:
    print( row )
```

Sample output:
```
{'Gateway': '0114A8C0', 'Mask': '00000000', 'Iface': 'wlan0', 'MTU': '0', 'Use': '0', 'IRTT': '0', 'RefCnt': '0', 'Metric': '303', 'Window': '0', 'Flags': '0003', 'Destination': '00000000'}
{'Gateway': '00000000', 'Mask': '000000FF', 'Iface': 'lo', 'MTU': '0', 'Use': '0', 'IRTT': '0', 'RefCnt': '0', 'Metric': '0', 'Window': '0', 'Flags': '0001', 'Destination': '0000007F'}
```
### php

Example:
```
<?php

$json_str = file_get_contents( "http://127.0.0.1:8080/route" );
$data_obj = json_decode( $json_str );
print_r( $data_obj );

?>
```

Sample output:
```
Array
(
    [0] => stdClass Object
        (
            [Destination] => 00000000
            [Flags] => 0003
            [Gateway] => 0114A8C0
            [IRTT] => 0
            [Iface] => wlan0
            [MTU] => 0
            [Mask] => 00000000
            [Metric] => 303
            [RefCnt] => 0
            [Use] => 0
            [Window] => 0
        )

    [1] => stdClass Object
        (
            [Destination] => 0000007F
            [Flags] => 0001
            [Gateway] => 00000000
            [IRTT] => 0
            [Iface] => lo
            [MTU] => 0
            [Mask] => 000000FF
            [Metric] => 0
            [RefCnt] => 0
            [Use] => 0
            [Window] => 0
        )

)
```


## See Also
