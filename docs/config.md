# Configuration

```yaml
scripts:
  sources:
    update_interval: 5s
    folder:
    <array>

datasources:
  clickhouse:
  <array>

  prometheus:
  <array>

channels:
  slack:
  <array>

global:
  send_start_notification:
    - slack1
  send_stop_notification:
    - slack1
```

## Global

- send_start_notification
- send_stop_notification

This options contains alert channels names for send message on start/stop balerter service

## Scripts Source

### Folder

|field|format|default|description|
|-|-|-|-|
|name|required, not empty|||
|path|required, not empty|||
|mask||*.lua||

#### example

```yaml
- name: 'source-name'
  path: '/path/to/folder'
  mask: '*.lua'
```

## Data Sources

### Clickhouse

|field|format|default|description|
|-|-|-|-|
|name|required, not empty|||
|host|required, not empty|||
|port|required, not zero|||
|username|required, not empty|||
|password||||
|database||||
|ssl_mode||||
|ssl_cert_path||||

#### example

```yaml
- name: ch1
  host: domain.com
  port: 6440
  username: user
  password: secret
  database: defaault
  ssl_mode: verified_full
  ssl_cert_path: /home/user/db.crt
```

### Prometheus

|field|format|default|description|
|-|-|-|-|
|name|required, not empty|||
|url|required, not empty|||
|basic_auth|||basic auth struct|

basic auth struct

|field|format|default|description|
|-|-|-|-|
|username||||
|password||||

#### example
```yaml
- name: prom1
	url: http://domain.com
	basic_auth:
		username: service_user
		password: QY2cvhcpCKPBnUwtPeNJUpkC
```

## Alert Channels

### Slack

|field|format|default|description|
|-|-|-|-|
|name|required, not empty|||
|url|required, not empty|||
|message_prefix_success|||will placed before success text in each message|
|message_prefix_error|||will placed before alert text in each message|
