local csv = {
    _VERSION = 'csv 1.2.0',
    _DESCRIPTION = 'CSV library for Lua',
    _URL         = 'https://github.com/FourierTransformer/ftcsv',
    _LICENSE     = [[
        The MIT License (MIT)
        Copyright (c) 2016-2020 Shakil Thakur
        Permission is hereby granted, free of charge, to any person obtaining a copy
        of this software and associated documentation files (the "Software"), to deal
        in the Software without restriction, including without limitation the rights
        to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
        copies of the Software, and to permit persons to whom the Software is
        furnished to do so, subject to the following conditions:
        The above copyright notice and this permission notice shall be included in all
        copies or substantial portions of the Software.
        THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
        IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
        FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
        AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
        LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
        OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
        SOFTWARE.
    ]]
}

-- perf
local sbyte = string.byte
local ssub = string.sub

-- luajit/lua compatability layer
local luaCompatibility = {}
if type(jit) == 'table' or _ENV then
    -- luajit and lua 5.2+
    luaCompatibility.load = _G.load
else
    -- lua 5.1
    luaCompatibility.load = loadstring
end

-- luajit specific speedups
-- luajit performs faster with iterating over string.byte,
-- whereas vanilla lua performs faster with string.find
if type(jit) == 'table' then
    luaCompatibility.LuaJIT = true
    -- finds the end of an escape sequence
    function luaCompatibility.findClosingQuote(i, inputLength, inputString, quote, doubleQuoteEscape)
        local currentChar, nextChar = sbyte(inputString, i), nil
        while i <= inputLength do
            nextChar = sbyte(inputString, i+1)

            -- this one deals with " double quotes that are escaped "" within single quotes "
            -- these should be turned into a single quote at the end of the field
            if currentChar == quote and nextChar == quote then
                doubleQuoteEscape = true
                i = i + 2
                currentChar = sbyte(inputString, i)

            -- identifies the escape toggle
            elseif currentChar == quote and nextChar ~= quote then
                return i-1, doubleQuoteEscape
            else
                i = i + 1
                currentChar = nextChar
            end
        end
    end

else
    luaCompatibility.LuaJIT = false

    -- vanilla lua closing quote finder
    function luaCompatibility.findClosingQuote(i, inputLength, inputString, quote, doubleQuoteEscape)
        local j, difference
        i, j = inputString:find('"+', i)
        if j == nil then
            return nil
        end
        difference = j - i
        if difference >= 1 then doubleQuoteEscape = true end
        if difference % 2 == 1 then
            return luaCompatibility.findClosingQuote(j+1, inputLength, inputString, quote, doubleQuoteEscape)
        end
        return j-1, doubleQuoteEscape
    end
end


-- determine the real headers as opposed to the header mapping
local function determineRealHeaders(headerField, fieldsToKeep) 
    local realHeaders = {}
    local headerSet = {}
    for i = 1, #headerField do
        if not headerSet[headerField[i]] then
            if fieldsToKeep ~= nil and fieldsToKeep[headerField[i]] then
                table.insert(realHeaders, headerField[i])
                headerSet[headerField[i]] = true
            elseif fieldsToKeep == nil then
                table.insert(realHeaders, headerField[i])
                headerSet[headerField[i]] = true
            end
        end
    end
    return realHeaders
end


local function determineTotalColumnCount(headerField, fieldsToKeep)
    local totalColumnCount = 0
    local headerFieldSet = {}
    for _, header in pairs(headerField) do
        -- count unique columns and
        -- also figure out if it's a field to keep
        if not headerFieldSet[header] and
            (fieldsToKeep == nil or fieldsToKeep[header]) then
            headerFieldSet[header] = true
            totalColumnCount = totalColumnCount + 1
        end
    end
    return totalColumnCount
end

local function generateHeadersMetamethod(finalHeaders)
    -- if a header field tries to escape, we will simply return nil
    -- the parser will still parse, but wont get the performance benefit of
    -- having headers predefined
    for _, headers in ipairs(finalHeaders) do
        if headers:find("]") then
            return nil
        end
    end
    local rawSetup = "local t, k, _ = ... \
    rawset(t, k, {[ [[%s]] ]=true})"
    rawSetup = rawSetup:format(table.concat(finalHeaders, "]] ]=true, [ [["))
    return luaCompatibility.load(rawSetup)
end

-- main function used to parse
local function parseString(inputString, i, options)

    -- keep track of my chars!
    local inputLength = options.inputLength or #inputString
    local currentChar, nextChar = sbyte(inputString, i), nil
    local skipChar = 0
    local field
    local fieldStart = i
    local fieldNum = 1
    local lineNum = 1
    local lineStart = i
    local doubleQuoteEscape, emptyIdentified = false, false

    local skipIndex
    local charPatternToSkip = "[" .. options.delimiter .. "\r\n]"

    --bytes
    local CR = sbyte("\r")
    local LF = sbyte("\n")
    local quote = sbyte('"')
    local delimiterByte = sbyte(options.delimiter)

    -- explode most used options
    local headersMetamethod = options.headersMetamethod
    local fieldsToKeep = options.fieldsToKeep
    local ignoreQuotes = options.ignoreQuotes
    local headerField = options.headerField
    local endOfFile = options.endOfFile
    local buffered = options.buffered

    local outResults = {}

    -- in the first run, the headers haven't been set yet.
    if headerField == nil then
        headerField = {}
        -- setup a metatable to simply return the key that's passed in
        local headerMeta = {__index = function(_, key) return key end}
        setmetatable(headerField, headerMeta)
    end

    if headersMetamethod then
        setmetatable(outResults, {__newindex = headersMetamethod})
    end
    outResults[1] = {}

    -- totalColumnCount based on unique headers and fieldsToKeep
    local totalColumnCount = options.totalColumnCount or determineTotalColumnCount(headerField, fieldsToKeep)

    local function assignValueToField()
        if fieldsToKeep == nil or fieldsToKeep[headerField[fieldNum]] then

            -- create new field
            if ignoreQuotes == false and sbyte(inputString, i-1) == quote then
                field = ssub(inputString, fieldStart, i-2)
            else
                field = ssub(inputString, fieldStart, i-1)
            end
            if doubleQuoteEscape then
                field = field:gsub('""', '"')
            end

            -- reset flags
            doubleQuoteEscape = false
            emptyIdentified = false

            -- assign field in output
            if headerField[fieldNum] ~= nil then
                outResults[lineNum][headerField[fieldNum]] = field
            else
                error('csv: too many columns in row ' .. options.rowOffset + lineNum)
            end
        end
    end

    while i <= inputLength do
        -- go by two chars at a time,
        --  currentChar is set at the bottom.
        nextChar = sbyte(inputString, i+1)

        -- empty string
        if ignoreQuotes == false and currentChar == quote and nextChar == quote then
            skipChar = 1
            fieldStart = i + 2
            emptyIdentified = true

        -- escape toggle.
        -- This can only happen if fields have quotes around them
        -- so the current "start" has to be where a quote character is.
        elseif ignoreQuotes == false and currentChar == quote and nextChar ~= quote and fieldStart == i then
            fieldStart = i + 1
            -- if an empty field was identified before assignment, it means
            -- that this is a quoted field that starts with escaped quotes
            -- ex: """a"""
            if emptyIdentified then
                fieldStart = fieldStart - 2
                emptyIdentified = false
            end
            skipChar = 1
            i, doubleQuoteEscape = luaCompatibility.findClosingQuote(i+1, inputLength, inputString, quote, doubleQuoteEscape)

        -- create some fields
        elseif currentChar == delimiterByte then
            assignValueToField()

            -- increaseFieldIndices
            fieldNum = fieldNum + 1
            fieldStart = i + 1

        -- newline
        elseif (currentChar == LF or currentChar == CR) then
            assignValueToField()

            -- handle CRLF
            if (currentChar == CR and nextChar == LF) then
                skipChar = 1
                fieldStart = fieldStart + 1
            end

            -- incrememnt for new line
            if fieldNum < totalColumnCount then
                -- sometimes in buffered mode, the buffer starts with a newline
                -- this skips the newline and lets the parsing continue.
                if buffered and lineNum == 1 and fieldNum == 1 and field == "" then
                    fieldStart = i + 1 + skipChar
                    lineStart = fieldStart
                else
                    error('csv: too few columns in row ' .. options.rowOffset + lineNum)
                end
            else
                lineNum = lineNum + 1
                outResults[lineNum] = {}
                fieldNum = 1
                fieldStart = i + 1 + skipChar
                lineStart = fieldStart
            end

        elseif luaCompatibility.LuaJIT == false then
            skipIndex = inputString:find(charPatternToSkip, i)
            if skipIndex then
                skipChar = skipIndex - i - 1
            end

        end

        -- in buffered mode and it can't find the closing quote
        -- it usually means in the middle of a buffer and need to backtrack
        if i == nil then
            if buffered then
                outResults[lineNum] = nil
                return outResults, lineStart
            else
                error("csv: can't find closing quote in row " .. options.rowOffset + lineNum ..
                 ". Try running with the option ignoreQuotes=true if the source incorrectly uses quotes.")
            end
        end

        -- Increment Counter
        i = i + 1 + skipChar
        if (skipChar > 0) then
            currentChar = sbyte(inputString, i)
        else
            currentChar = nextChar
        end
        skipChar = 0
    end

    if buffered and not endOfFile then
        outResults[lineNum] = nil
        return outResults, lineStart
    end

    -- create last new field
    assignValueToField()

    -- remove last field if empty
    if fieldNum < totalColumnCount then

        -- indicates last field was really just a CRLF,
        -- so, it can be removed
        if fieldNum == 1 and field == "" then
            outResults[lineNum] = nil
        else
            error('csv: too few columns in row ' .. options.rowOffset + lineNum)
        end
    end

    return outResults, i, totalColumnCount
end

local function handleHeaders(headerField, options)
    -- make sure a header isn't empty
    for _, headerName in ipairs(headerField) do
        if #headerName == 0 then
            error('csv: Cannot parse a file which contains empty headers')
        end
    end

    -- for files where there aren't headers!
    if options.headers == false then
        for j = 1, #headerField do
            headerField[j] = j
        end
    end

    -- rename fields as needed!
    if options.rename then
        -- basic rename (["a" = "apple"])
        for j = 1, #headerField do
            if options.rename[headerField[j]] then
                headerField[j] = options.rename[headerField[j]]
            end
        end
        -- files without headers, but with a options.rename need to be handled too!
        if #options.rename > 0 then
            for j = 1, #options.rename do
                headerField[j] = options.rename[j]
            end
        end
    end

    -- apply some sweet header manipulation
    if options.headerFunc then
        for j = 1, #headerField do
            headerField[j] = options.headerFunc(headerField[j])
        end
    end

    return headerField
end

-- load an entire file into memory
local function loadFile(textFile, amount)
    local file = io.open(textFile, "r")
    if not file then error("csv: File not found at " .. textFile) end
    local lines = file:read(amount)
    if amount == "*all" then
        file:close()
    end
    return lines, file
end

local function initializeInputFromStringOrFile(inputFile, options, amount)
    -- handle input via string or file!
    local inputString, file
    if options.loadFromString then
        inputString = inputFile
    else
        inputString, file = loadFile(inputFile, amount)
    end

    -- if they sent in an empty file...
    if inputString == "" then
        error('csv: Cannot parse an empty file')
    end
    return inputString, file
end

local function parseOptions(delimiter, options, fromParseLine)
    -- delimiter MUST be one character
    assert(#delimiter == 1 and type(delimiter) == "string", "the delimiter must be of string type and exactly one character")

    local fieldsToKeep = nil

    if options then
        if options.headers ~= nil then
            assert(type(options.headers) == "boolean", "csv only takes the boolean 'true' or 'false' for the optional parameter 'headers' (default 'true'). You passed in '" .. tostring(options.headers) .. "' of type '" .. type(options.headers) .. "'.")
        end
        if options.rename ~= nil then
            assert(type(options.rename) == "table", "csv only takes in a key-value table for the optional parameter 'rename'. You passed in '" .. tostring(options.rename) .. "' of type '" .. type(options.rename) .. "'.")
        end
        if options.fieldsToKeep ~= nil then
            assert(type(options.fieldsToKeep) == "table", "csv only takes in a list (as a table) for the optional parameter 'fieldsToKeep'. You passed in '" .. tostring(options.fieldsToKeep) .. "' of type '" .. type(options.fieldsToKeep) .. "'.")
            local ofieldsToKeep = options.fieldsToKeep
            if ofieldsToKeep ~= nil then
                fieldsToKeep = {}
                for j = 1, #ofieldsToKeep do
                    fieldsToKeep[ofieldsToKeep[j]] = true
                end
            end
            if options.headers == false and options.rename == nil then
                error("csv: fieldsToKeep only works with header-less files when using the 'rename' functionality")
            end
        end
        if options.loadFromString ~= nil then
            assert(type(options.loadFromString) == "boolean", "csv only takes a boolean value for optional parameter 'loadFromString'. You passed in '" .. tostring(options.loadFromString) .. "' of type '" .. type(options.loadFromString) .. "'.")
        end
        if options.headerFunc ~= nil then
            assert(type(options.headerFunc) == "function", "csv only takes a function value for optional parameter 'headerFunc'. You passed in '" .. tostring(options.headerFunc) .. "' of type '" .. type(options.headerFunc) .. "'.")
        end
        if options.ignoreQuotes == nil then
            options.ignoreQuotes = false
        else
            assert(type(options.ignoreQuotes) == "boolean", "csv only takes a boolean value for optional parameter 'ignoreQuotes'. You passed in '" .. tostring(options.ignoreQuotes) .. "' of type '" .. type(options.ignoreQuotes) .. "'.")
        end
        if options.bufferSize ~= nil then
            assert(type(options.bufferSize) == "number", "csv only takes a number value for optional parameter 'bufferSize'. You passed in '" .. tostring(options.bufferSize) .. "' of type '" .. type(options.bufferSize) .. "'.")
            if fromParseLine == false then
                error("csv: bufferSize can only be specified using 'parseLine'. When using 'parse', the entire file is read into memory")
            end
        end
    else
        options = {
            ["headers"] = true,
            ["loadFromString"] = false,
            ["ignoreQuotes"] = false,
            ["bufferSize"] = 2^16
        }
    end

    return options, fieldsToKeep

end

local function findEndOfHeaders(str, entireFile)
    local i = 1
    local quote = sbyte('"')
    local newlines = {
        [sbyte("\n")] = true,
        [sbyte("\r")] = true
    }
    local quoted = false
    local char = sbyte(str, i)
    repeat
        -- this should still work for escaped quotes
        -- ex: " a "" b \r\n " -- there is always a pair around the newline
        if char == quote then
            quoted = not quoted
        end
        i = i + 1
        char = sbyte(str, i)
    until (newlines[char] and not quoted) or char == nil

    if not entireFile and char == nil then
        error("csv: bufferSize needs to be larger to parse this file")
    end

    local nextChar = sbyte(str, i+1)
    if nextChar == sbyte("\n") and char == sbyte("\r") then
        i = i + 1
    end
    return i
end

local function determineBOMOffset(inputString)
    -- BOM files start with bytes 239, 187, 191
    if sbyte(inputString, 1) == 239
        and sbyte(inputString, 2) == 187
        and sbyte(inputString, 3) == 191 then
        return 4
    else
        return 1
    end
end

local function parseHeadersAndSetupArgs(inputString, delimiter, options, fieldsToKeep, entireFile)
    local startLine = determineBOMOffset(inputString)

    local endOfHeaderRow = findEndOfHeaders(inputString, entireFile)

    local parserArgs = {
        delimiter = delimiter,
        headerField = nil,
        fieldsToKeep = nil,
        inputLength = endOfHeaderRow,
        buffered = false,
        ignoreQuotes = options.ignoreQuotes,
        rowOffset = 0
    }

    local rawHeaders, endOfHeaders = parseString(inputString, startLine, parserArgs)

    -- manipulate the headers as per the options
    local modifiedHeaders = handleHeaders(rawHeaders[1], options)
    parserArgs.headerField = modifiedHeaders
    parserArgs.fieldsToKeep = fieldsToKeep
    parserArgs.inputLength = nil

    if options.headers == false then endOfHeaders = startLine end

    local finalHeaders = determineRealHeaders(modifiedHeaders, fieldsToKeep)
    if options.headers ~= false then
        local headersMetamethod = generateHeadersMetamethod(finalHeaders)
        parserArgs.headersMetamethod = headersMetamethod
    end

    return endOfHeaders, parserArgs, finalHeaders
end

-- runs the show!
function csv.parse(inputFile, delimiter, options)
    local options, fieldsToKeep = parseOptions(delimiter, options, false)

    local inputString = initializeInputFromStringOrFile(inputFile, options, "*all")

    local endOfHeaders, parserArgs, finalHeaders = parseHeadersAndSetupArgs(inputString, delimiter, options, fieldsToKeep, true)

    local output = parseString(inputString, endOfHeaders, parserArgs)

    return output, finalHeaders
end

local function getFileSize (file)
    local current = file:seek()
    local size = file:seek("end")
    file:seek("set", current)
    return size
end

local function determineAtEndOfFile(file, fileSize)
    if file:seek() >= fileSize then
        return true
    else
        return false
    end
end

local function initializeInputFile(inputString, options)
    if options.loadFromString == true then
        error("csv: parseLine currently doesn't support loading from string")
    end
    return initializeInputFromStringOrFile(inputString, options, options.bufferSize)
end

function csv.parseLine(inputFile, delimiter, userOptions)
    local options, fieldsToKeep = parseOptions(delimiter, userOptions, true)
    local inputString, file = initializeInputFile(inputFile, options)


    local fileSize, atEndOfFile = 0, false
    fileSize = getFileSize(file)
    atEndOfFile = determineAtEndOfFile(file, fileSize)

    local endOfHeaders, parserArgs, _ = parseHeadersAndSetupArgs(inputString, delimiter, options, fieldsToKeep, atEndOfFile)
    parserArgs.buffered = true
    parserArgs.endOfFile = atEndOfFile

    local parsedBuffer, endOfParsedInput, totalColumnCount = parseString(inputString, endOfHeaders, parserArgs)
    parserArgs.totalColumnCount = totalColumnCount

    inputString = ssub(inputString, endOfParsedInput)
    local bufferIndex, returnedRowsCount = 0, 0
    local currentRow, buffer

    return function()
        -- check parsed buffer for value
        bufferIndex = bufferIndex + 1
        currentRow = parsedBuffer[bufferIndex]
        if currentRow then
            returnedRowsCount = returnedRowsCount + 1
            return returnedRowsCount, currentRow
        end

        -- read more of the input
        buffer = file:read(options.bufferSize)
        if not buffer then
            file:close()
            return nil
        else
            parserArgs.endOfFile = determineAtEndOfFile(file, fileSize)
        end

        -- appends the new input to what was left over
        inputString = inputString .. buffer

        -- re-analyze and load buffer
        parserArgs.rowOffset = returnedRowsCount
        parsedBuffer, endOfParsedInput = parseString(inputString, 1, parserArgs)
        bufferIndex = 1

        -- cut the input string down
        inputString = ssub(inputString, endOfParsedInput)

        if #parsedBuffer == 0 then
            error("csv: bufferSize needs to be larger to parse this file")
        end

        returnedRowsCount = returnedRowsCount + 1
        return returnedRowsCount, parsedBuffer[bufferIndex]
    end
end



-- The ENCODER code is below here
-- This could be broken out, but is kept here for portability


local function delimitField(field)
    field = tostring(field)
    if field:find('"') then
        return field:gsub('"', '""')
    else
        return field
    end
end

local function escapeHeadersForLuaGenerator(headers)
    local escapedHeaders = {}
    for i = 1, #headers do
        if headers[i]:find('"') then
            escapedHeaders[i] = headers[i]:gsub('"', '\\"')
        else
            escapedHeaders[i] = headers[i]
        end
    end
    return escapedHeaders
end

-- a function that compiles some lua code to quickly print out the csv
local function csvLineGenerator(inputTable, delimiter, headers)
    local escapedHeaders = escapeHeadersForLuaGenerator(headers)

    local outputFunc = [[
        local args, i = ...
        i = i + 1;
        if i > ]] .. #inputTable .. [[ then return nil end;
        return i, '"' .. args.delimitField(args.t[i]["]] ..
            table.concat(escapedHeaders, [["]) .. '"]] ..
            delimiter .. [["' .. args.delimitField(args.t[i]["]]) ..
            [["]) .. '"\r\n']]

    local arguments = {}
    arguments.t = inputTable
    -- we want to use the same delimitField throughout,
    -- so we're just going to pass it in
    arguments.delimitField = delimitField

    return luaCompatibility.load(outputFunc), arguments, 0

end

local function validateHeaders(headers, inputTable)
    for i = 1, #headers do
        if inputTable[1][headers[i]] == nil then
            error("csv: the field '" .. headers[i] .. "' doesn't exist in the inputTable")
        end
    end
end

local function initializeOutputWithEscapedHeaders(escapedHeaders, delimiter)
    local output = {}
    output[1] = '"' .. table.concat(escapedHeaders, '"' .. delimiter .. '"') .. '"\r\n'
    return output
end

local function escapeHeadersForOutput(headers)
    local escapedHeaders = {}
    for i = 1, #headers do
        escapedHeaders[i] = delimitField(headers[i])
    end
    return escapedHeaders
end

local function extractHeadersFromTable(inputTable)
    local headers = {}
    for key, _ in pairs(inputTable[1]) do
        headers[#headers+1] = key
    end

    -- lets make the headers alphabetical
    table.sort(headers)

    return headers
end

local function getHeadersFromOptions(options)
    local headers = nil
    if options then
        if options.fieldsToKeep ~= nil then
            assert(
                type(options.fieldsToKeep) == "table", "csv only takes in a list (as a table) for the optional parameter 'fieldsToKeep'. You passed in '" .. tostring(options.headers) .. "' of type '" .. type(options.headers) .. "'.")
            headers = options.fieldsToKeep
        end
    end
    return headers
end

local function initializeGenerator(inputTable, delimiter, options)
    -- delimiter MUST be one character
    assert(#delimiter == 1 and type(delimiter) == "string", "the delimiter must be of string type and exactly one character")

    local headers = getHeadersFromOptions(options)
    if headers == nil then
        headers = extractHeadersFromTable(inputTable)
    end
    validateHeaders(headers, inputTable)

    local escapedHeaders = escapeHeadersForOutput(headers)
    local output = initializeOutputWithEscapedHeaders(escapedHeaders, delimiter)
    return output, headers
end

-- works really quickly with luajit-2.1, because table.concat life
function csv.encode(inputTable, delimiter, options)
    local output, headers = initializeGenerator(inputTable, delimiter, options)

    for i, line in csvLineGenerator(inputTable, delimiter, headers) do
        output[i+1] = line
    end

    -- combine and return final string
    return table.concat(output)
end

return csv
