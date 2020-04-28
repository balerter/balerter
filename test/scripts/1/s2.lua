-- @interval 5s

local db = require('datasource.clickhouse.ch1')
local log = require('log')
local csv = require('csv')
local h = require('h')

res, err = db.query('SELECT * FROM users')
if err ~= nil then
    log.error('query error: ' .. err)
    return
end

h.print(csv.encode(res, ","))