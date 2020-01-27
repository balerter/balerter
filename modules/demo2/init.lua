-- External lua module in folder
-- Initial file should be init.lua
--
-- Usage:
--
-- local demo2 = require('demo2')
-- local res = demo.bar()
-- res will be equals 'baz'

local M = {}

local function bar()
    return 'baz'
end

rawset(M, 'bar', bar)

return M