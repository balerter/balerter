# Core module`KV`

**KV** - it's key/value storage.

You can use this module for share some values between different scripts or store some data for a single script between executions.

Usage:
```
local kv = require('kv')

local value, err = kv.get('some string key')

err = kv.put('some string key', 'some value')
```

## API

### `get(<KEY>) value, error`

Get value by key

```
local kv = require('kv')
local val, err = kv.get('key1')
```

Returns an error if the variable does not exists

### `put(<KEY>, <VALUE>) error`

Put a new variable to storage

```
local kv = require('kv')
local err1 = kv.put('key1', 'value')
local err2 = kv.put('key2', 42)
```

Returns an error if the variable already exists

### `upsert(<KEY>, <VALUE>) error`

Put a variable to storage. The variable will be rewrite if exists

```
local kv = require('kv')
local err = kv.upsert('key2', 'value')
```

### `delete(<KEY>) error`

Delete a variable from storage

```
local kv = require('kv')
local err = kv.delete('key2')
```

Returns an error if the variable does not exists
