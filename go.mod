module goormbenchorm

go 1.15

// replace github.com/webx-top/db => ../../webx-top/db

require (
	github.com/beego/beego/v2 v2.0.2
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gocraft/dbr v0.0.0-20190714181702-8114670a83bd
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.6
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/upper/db/v4 v4.6.0
	github.com/webx-top/db v1.23.8
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gorm.io/driver/mysql v1.3.3
	gorm.io/driver/postgres v1.3.4
	gorm.io/gorm v1.23.4
	mellium.im/sasl v0.2.1 // indirect
	xorm.io/builder v0.3.10 // indirect
	xorm.io/xorm v1.3.0
)
