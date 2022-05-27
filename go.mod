module go-web-kabu-gin-gee

go 1.18
require kabucache v0.0.0
replace kabucache   => ./kabucache
require lru v0.0.0
replace lru    => ./kabucache/lru