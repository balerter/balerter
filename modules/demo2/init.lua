local M = {}

local function bar()
    return 'baz'
end

rawset(M, 'bar', bar)

return M