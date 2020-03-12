# balerter [![Go Report Card](https://goreportcard.com/badge/github.com/balerter/balerter)](https://goreportcard.com/report/github.com/balerter/balerter) ![Test](https://github.com/balerter/balerter/workflows/Test/badge.svg)

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
    -v /peth/to/config.yml:/opt/config.yml \
    -v /path/to/scripts:/opt/scripts \ 
    -v /path/to/cert.crt:/home/user/db.crt \
    balerter/balerter -config=/opt/config.yml
```

Config file `config.yml`
```yaml
scripts:
  sources:
    update_interval: 5s
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
      ssl_mode: verified_full
      ssl_cert_path: /home/user/db.crt

channels:
  slack:
    - name: slack1
      url: https://hooks.slack.com/services/hash
      message_prefix_success: ':eight_spoked_asterisk: '
      message_prefix_error: ':sos: '
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
    alert.error("rps-min-limit", "Requests RPS are very small: " .. tostring(resultRPS))
else
    alert.success("rps-min-limit", "Requests RPS ok"")
end 
```

## Roadmap

### For v1.0.0

- [x] add filters to api `/alerts`
    - [x] by alert level
    - [x] by alert name 
- [x] stabilize core DB modules
- [ ] full translate the documentation to English
- [x] add core module `http` for send requests from scripts

### Other

- core modules enhancements
    - [ ] add prometheus methods for querying metadata
        - [ ] series
        - [ ] labels
        - [ ] label values
- [ ] Official grafana dashboard
- New script sources
    - [ ] Postgres - select scripts from Postgres table
    - [ ] ...
- new entity: external KV storage - for persist data
- engines for external KV storage
    - [ ] file
    - [ ] consul
    - [ ] postgres
- New datasources
    - [ ] MongoDB
    - [ ] ...
- New channels
    - [ ] email
    - [ ] webhook
    - [ ] rocketchat
    - [ ] ...
- [ ] chart module enhancement