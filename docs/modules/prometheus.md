# Interval module Prometheus

## Usage

```lua
local prom = require('datasource.prometheus.<DATASOURCE_NAME')
```

## Methods

### Query

```lua
local res, err = prom.query('<PromQL query')
```

Return values:
- err - error string or `nil` if no error
- res - result table

> todo example

### Range

```lua
local res, err = prom.range('<PromQL query')
```

Return values:
- err - error string or `nil` if no error
- res - result table

> todo example

