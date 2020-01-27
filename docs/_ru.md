# Balerter - a script based alerting system

Для описания правил срабатывания оповещений испльзуются Lua-скрипты.

Это позволяет:
- получить данные одновременно из разничных источников
- написать сложную бизнес-логику 
- гибко настроить уведомеления. Например, вы отслеживаете показатель запросов в секунду. И вы можете настроить, что при падении значения до уровня X уведомеление отправляется дежурному инженеру, а при достижении уровня Y - еще и к руководителю 

#### Дополнительные модули

Если у вас есть общая логика, которую вы хотите переиспользовать в различных скриптах, вы можете написать свои Lua-модули и поместить их в /modules

После этого, вы можете импортировать и использовать свои модули в скриптах. В данной папке есть примеры таких модулей 

#### Пример использования:

Получаем некоторую статистику по запросам из Prometheus.
Для тех показателей, которые ниже порогового значения, получаем дополнительные данные из Postgres (название клиента) и отправляем оповещение.
С помощью мета-тега `@interval` указываем, что скрипт исполняется каждую минуту.

```
-- @interval 1m
-- @name Client Requests Level

local prom1 = require('datasource.prometheus.prom1')
local pg1 = require('datasource.postgres.pg1')
local alert = require('alert')
local log = require('log')

local promData, err = prom1.query('sum(requests) by (client_id)')
if err ~= nil then
    log.error('error quering prom1: ' .. err)
    return
end

local clientsData, err = pg1.query('SELECT id, client_name FROM clients')
if err ~= nil then
    log.error('error quering pg1: ' .. err)
    return
end

-- ... convert clientsData to table {'1' => 'client1 name', '2' => 'client2 name' ...}

local min_requests = 100

for _, item in paris(prom_data) do
    local alertName = 'requests-for-client-' .. tostring(item.value)

    if item.value < min_requests then
        alert.on(alertName, 'Client ' .. clientData[item.labels.client_id] .. ' has low requests level')
    else
        alert.off(alertName, 'Client ' .. clientData[item.labels.client_id] .. ' has normal requests level')
    end 
end
```

Оповещение будет отправляться, когда наше условие будет меняться. То есть, после первого срабатывания условия, когда значение меньше 100, будет отправлено оповещение о включении алерта. Оно не будет повторяться, пока значение не поднимется выше 100. В этом случае будет отправлено оповещение о выключении алерта

Имеется возможность более гибко настраивать алерты, если это необходимо.

Например, для критичных алертов, можно указать, чтобы они отправлялись каждый раз (или один раз за каждые N срабатываний) 


## Источники данных

Реализовано:

- `prometheus`
- `clickhouse`
- `postgres`

Рассматриваются к реализации:
- `http api`
- `mysql`

## Источники скриптов

Реализовано:

- `filesystem folder` - можно указать маску для выбора определенных файлов 

Рассматриваются к реализации:

- `postgres`
- `mysql`
- `http api`

## Каналы оповещений

Реализовано:

- `slack`
 
Рассматриваются к реализации:

- `email`
- `telegram`
- `http api`


