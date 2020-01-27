-- External lua module in single file
--
-- Usage:
--
-- local demo = require('demo')
-- local res = demo.foo()
-- res will be equals 'bar'

local M = {}

local function foo()
    return 'bar'
end

rawset(M, 'foo', foo)

return M