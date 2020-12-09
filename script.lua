local lua2go = require('lua2go')
local sign = lua2go.Load('./sign.so')

lua2go.Externs [[
extern char* Hello();
extern char* Sign(GoString p0);
]]

local charset = {} do -- [0-9a-zA-Z]
    for c = 48, 57 do table.insert(charset, string.char(c)) end
    for c = 65, 90 do table.insert(charset, string.char(c)) end
    for c = 97, 122 do table.insert(charset, string.char(c)) end
end

local function randomString(length)
    if not length or length <= 0 then return '' end
    math.randomseed(os.clock() ^ 5)
    return randomString(length - 1) .. charset[math.random(1, #charset)]
end

function request()
    local header = {}
    local randomData = randomString(32)
    local goSign = sign.Sign(lua2go.ToGoString(randomData))
    lua2go.AddToGC(goSign)
    header["signature"] = randomData .. ';' .. lua2go.ToLua(goSign)
    return wrk.format("POST", wrk.path, header)
end
