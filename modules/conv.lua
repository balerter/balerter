local M = {}

-- Convert table to a map with keyField value as a key
-- Use for database results
--
-- returns: result table and error if occurred
-- example:
--
-- local t = { { 'id' = 1, 'name' = 'foo' }, { 'id' = 2, 'name' = 'bar' } }
-- result = conv.toMap(t, 'id')
-- result: { '1' = { 'id' = 1, 'name' = 'foo' }, '2' = { 'id' = 2, 'name' = 'bar' } }
local function toMap(t, keyField)
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

rawset(M, 'toMap', toMap)

return M