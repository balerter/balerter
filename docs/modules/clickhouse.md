# Internal module Clickhouse

## Usage

```lua
local ch = require('datasource.clickhouse.<datasource_name>')
```

## Methods

### Query

```lua
local res, err = ch.query('<SQL query')
```

Return values:
- err - error string or `nil` if no error
- res - result table

```
{
    {
        '<FIELD_NAME>' = <FIELD_VALUE>,
        ...
    },
    ...
}
``` 

For example, query `local res, err = ch.query("SELECT table, sum(bytes) AS size FROM system.parts WHERE active AND database = 'system' GROUP BY table")` will returns `res`:  
```
{
    {
        table = 'query_thread_log',
        size = 5478131511
    },
    {
        table = 'metric_log',
        size = 486005337
    },
    {
        table = 'query_log',
        size = 3867910587
    },
    {
        table = 'part_log',
        size = 3383019302
    }    
}

```