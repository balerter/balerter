-- @interval 5s

local pg = require('datasource.postgres.pg1')
local log = require('log')
local h = require('h')

res, err = pg.query('SELECT * FROM users')
if err ~= nil then
    log.error('query error: ' .. err)
    return
end

h.printTable(res)