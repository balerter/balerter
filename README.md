![GitHub release (latest by date)](https://img.shields.io/github/v/release/balerter/balerter) [![Go Report Card](https://goreportcard.com/badge/github.com/balerter/balerter)](https://goreportcard.com/report/github.com/balerter/balerter) ![Test](https://github.com/balerter/balerter/workflows/Test/badge.svg) [![codecov](https://codecov.io/gh/balerter/balerter/branch/master/graph/badge.svg)](https://codecov.io/gh/balerter/balerter) 

![logo.png](logo.png)

> A Project in active development. Features may have breaking changes at any time before v1.0.0 version 

- [Telegram Group](https://t.me/balerter)

Balerter is a scripts based alerting system.

In your script you may:
- obtain needed data from different data sources (prometheus, clickhouse, postgres, external HTTP API etc.)
- analyze data and make a decision about alert status
- change Alerts statuses and receive notifications about it 

In the example bellow we create one Clickhouse datasource, one scripts source and one alert channel.
In the script we run query to clickhouse, check the value and fire the alert (or switch off it)   

Full documentation available on https://balerter.com

## Example

```
docker pull balerter/balerter
```

```
docker run \
    -v /path/to/config.yml:/opt/config.yml \
    -v /path/to/scripts:/opt/scripts \ 
    -v /path/to/cert.crt:/home/user/db.crt \
    balerter/balerter -config=/opt/config.yml
```

Config file `config.yml`
```yaml
scripts:
  sources:
    updateInterval: 5s
    folder:
      - name: debug-folder
        path: /opt/scripts
        mask: '*.lua'

datasources:
  clickhouse:
    - name: ch1
      host: localhost
      port: 6440
      username: default
      password: secret
      database: default
      sslMode: verified_full
      sslCertPath: /home/user/db.crt

channels:
  slack:
    - name: slack1
      url: https://hooks.slack.com/services/hash
```

Sample script `rps.lua`
```
-- @cront */10 * * * *
-- @name script1

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
    alert.error("rps-min-limit", "Requests RPS are very small: " .. tostring(resultRPS))
else
    alert.success("rps-min-limit", "Requests RPS ok")
end 
```

Also, you can to write tests!

An example:

```
-- @test script1
-- @name script1-test

test = require('test')

local resp = {
    {
        rps = 10
    }
} 

test.datasource('clickhouse.ch1').on('query', 'SELECT sum(requests) AS rps FROM some_table WHERE date = now()').response(resp)

test.alert().assertCalled('error', 'rps-min-limit', 'Requests RPS are very small: 10')
test.alert().assertNotCalled('success', 'rps-min-limit', 'Requests RPS ok')
```

See a documentation on [https://balerter.com](https://balerter.com)
