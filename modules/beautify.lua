-- Pretty print tables
--
-- Usage:
-- beautify.table(variable)
--

local M = {}

local paddingWith = '    '

local function _table(data, tab)
    local result = ''

    result = result .. '{\n'

    for key, value in pairs(data) do
        result = result .. string.rep(paddingWith, tab) .. tostring(key) .. ' = '
        local r = ''
        if type(value) == 'table' then
            r = _table(value, tab + 1)
        else
            r = tostring(value) .. '\n'
        end
        result = result .. r
    end

    result = result .. string.rep(paddingWith, tab - 1) .. '}\n'

    return result
end

local function exportTable(data)
    if type(data) ~= 'table' then
        error('wrong data type')
        return
    end
    print(_table(data, 1))
end


rawset(M, 'table', exportTable)

return M
