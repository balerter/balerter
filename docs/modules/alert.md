# Internal module Alert

Alert module allows switch on/off alerts. 

## Usage

```lua
local alert = require('alert')
```

## Methods

### on

```lua
alert.on('<UNIQUE ALERT NAME>', '<ALERT MESSAGE>')
```

Switching on alert. If alert had status 'off', will be send `error` notification message

### off

```lua
alert.off('<UNIQUE ALERT NAME>', '<ALERT MESSAGE>')
```

Switching off alert. If alert had status 'on', will be send `success` notification message
