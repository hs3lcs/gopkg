package dbms

type Cfg struct {
	// db:read-write
	MYSQL_RW_HOST string
	MYSQL_RW_USER string
	MYSQL_RW_PASS string
	// db:read-only
	MYSQL_RO_HOST string
	MYSQL_RO_USER string
	MYSQL_RO_PASS string
	// redis cache
	REDIS_HOST string
	REDIS_PASS string
}

var Config *Cfg
