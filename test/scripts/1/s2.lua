-- @interval 5s

local alert = require('alert')
local db = require('datasource.clickhouse.ch1')
local log = require('log')
local h = require('h')

log.info('clickhouse query: SELECT * FROM users')

res, err = db.query('SELECT * FROM users')
if err ~= nil then
    log.error('query error: ' .. err)
    return
end

alert.warn('alert-id', 'An test warning from clickhouse query occured!')
alert.success('alert-id', 'Test warning from clickhouse query is gone!')
h.print(res)