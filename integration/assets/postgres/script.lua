local pg1 = require('datasource.postgres.pg1')
local log = require('log')
local h = require('h')

res, err = pg1.query('SELECT * FROM users')
if err ~= nil then
    log.error('query error: ' + err)
    return
end

h.print(res)