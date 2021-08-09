wrk.method="GET"

local queries = {
    "/getBestStoresList?longitude=106&latitude=-6&listNum=20&method=3",
   "/getBestStoresList?longitude=110&latitude=6&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=-8&listNum=10&method=3",
   "/getBestStoresList?longitude=106&latitude=-7&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=-6&listNum=30&method=3",
   "/getBestStoresList?longitude=106&latitude=-6&listNum=10&method=3",
   "/getBestStoresList?longitude=106&latitude=-4&listNum=20&method=3",
   "/getBestStoresList?longitude=200&latitude=-6&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=-1&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=0&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=-6&listNum=30&method=3",
   "/getBestStoresList?longitude=106&latitude=-6&listNum=10&method=3",
   "/getBestStoresList?longitude=106&latitude=-6&listNum=20&method=3",
   "/getBestStoresList?longitude=106&latitude=-8&listNum=20&method=3",
}

local i = 0

request = function()
    local path = wrk.format(nil, queries[i % #queries + 1])
    i = i + 1
    return path
end
