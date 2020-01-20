local M = {}

local function foo()
    return 'bar'
end

rawset(M, 'foo', foo)

return M