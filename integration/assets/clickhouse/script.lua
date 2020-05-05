local ch1 = require('datasource.clickhouse.ch1')
local log = require('log')
local h = require('h')

res, err = ch1.query('SELECT * FROM users ORDER BY id')
if err ~= nil then
    log.error('query error: ' + err)
    return
end

h.print(res)