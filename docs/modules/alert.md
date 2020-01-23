# Internal module Alert

Alert module allows switch on/off alerts.

A Message has a alert name and a message text.

Alert name must be unique. 

A message  may also contain Fields - an array of string

## Usage

```lua
local alert = require('alert')
```

## Methods

### on

```
alert.on('<UNIQUE ALERT NAME>', '<ALERT MESSAGE>'[, <ALERT FIELDS TABLE>])
```

Switching on alert. If alert had status 'off', will be send `error` notification message

### off

```
alert.off('<UNIQUE ALERT NAME>', '<ALERT MESSAGE>'[, <ALERT FIELDS TABLE>])
```

Switching off alert. If alert had status 'on', will be send `success` notification message

Examples:

```
alert.on('cpu-90', 'CPU busy more than 90%')

alert.on('mem', 'Memory data', {'node1', 'node2'} )

alert.off('cpu-90', 'CPU is ok')
```