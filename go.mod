module github.com/chremoas/migration-tool

go 1.15

require (
	github.com/chremoas/auth-srv v1.3.1
	github.com/doug-martin/goqu/v9 v9.10.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.10
	github.com/jmoiron/sqlx v1.2.0
)

replace github.com/chremoas/auth-srv => ../auth-srv
