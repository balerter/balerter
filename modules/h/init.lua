local M = {}

local function _printTable(data, tab)
    local paddingWith = '    '
    local result = ''

    result = result .. '{\n'

    for key, value in pairs(data) do
        result = result .. string.rep(paddingWith, tab) .. tostring(key) .. ' = '
        local r = ''
        if type(value) == 'table' then
            r = _printTable(value, tab + 1)
        else
            r = tostring(value) .. '\n'
        end
        result = result .. r
    end

    result = result .. string.rep(paddingWith, tab - 1) .. '}\n'

    return result
end

local function printTable(data)
    if type(data) ~= 'table' then
        error('wrong data type')
        return
    end
    print(_printTable(data, 1))
end

local function tableToMap(t, keyField)
    local result = {}

    for _, item in pairs(t) do
        local key = item[keyField]
        if key == nil then
            return nil, error('key field not found')
        end

        result[key] = item
    end

    return result
end

-- tableToMap(t, keyField)
--
-- Convert a table to a map with 'keyField' value as a key
-- Use for database results
--
-- Returns: a result table and an error if occurred
-- An example:
--
-- local t = { { 'id' = 1, 'name' = 'foo' }, { 'id' = 2, 'name' = 'bar' } }
-- result = tableToMap(t, 'id')
-- result: { '1' = { 'id' = 1, 'name' = 'foo' }, '2' = { 'id' = 2, 'name' = 'bar' } }
rawset(M, 'tableToMap', tableToMap)

-- Pretty print tables
--
-- Usage:
-- printTable(variable)
rawset(M, 'printTable', printTable)

return M