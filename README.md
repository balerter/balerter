# balerter  

Balerter is a scripts based alerting system.

> todo 

## Example

Config file `config.yml`
```yaml
scripts:
  sources:
    update_interval: 5s
    folder:
      - name: debug-folder
        path: /home/user/scripts
        mask: '*.lua'

datasources:
  clickhouse:
    - name: ch1
      host: localhost
      port: 6440
      username: default
      password: secret
      database: default
      ssl_mode: verified_full
      ssl_cert_path: /home/user/db.crt

channels:
  slack:
    - name: slack1
      url: https://hooks.slack.com/services/hash
      message_prefix_success: ':eight_spoked_asterisk: '
      message_prefix_error: ':sos: '

global:
  send_start_notification:
    - slack1
  send_stop_notification:
    - slack1
```

Sample script `rps.lua`
```
-- @interval 10s

local minRequestsRPS = 100

local log = require("log")
local ch1 = require("datasource.clickhouse.ch1")

local res, err = ch1.query("SELECT sum(requests) AS rps FROM some_table WHERE date = now()")
if err ~= nil then
    log.error("clickhouse 'ch1' query error: " .. err)
    return
end

local resultRPS = res[1].rps

if resultRPS < minResultRPS then
    alert.on("rps-min-limit", "Requests RPS are very small: " .. tostring(resultRPS))
else
    alert.off("rps-min-limit", "Requests RPS ok"")
end 
```

## Modules

Internal modules are divided into three types:
- Data Source
- Scripts Source
- Alert Channel

**Data Source** allows give access to data for analyze 

**Script Source** allows obtain scripts for run

**Alert Channel** allows send notifications  

### internal modules support

Currently supports:

|data source|script source|alert channel|
|-----------|-------------|-------------|
| clickhouse |filesystem folder |slack |
| prometheus| | |​

Plans to support:

|data source|script source|alert channel|
|-----------|-------------|-------------|
|postgres| |email|
|http| |telegram|

Possible plans to support:

|data source|script source|alert channel|
|-----------|-------------|-------------|
|mysql|postgres|webpush|
| | |whatsapp|

### external modules

Also supports external LUA-script modules. You should place it into `./modules` folder.
In `./modules` folder present two demo modules: demo and demo2.
You can use it by follow example:
```
local demo = require('demo')
local demo2 = require('demo2')

print(demo.foo())
print(demo2.bar())
```

You can place into this folder your own modules and use it. 
More modules can be found in the repo https://github.com/balerter/modules

## Documentation

- [configuration](docs/config.md)