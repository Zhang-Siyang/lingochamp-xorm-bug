module xorm_max_id_bug

go 1.22.6

require (
	github.com/lingochamp/core v0.5.8-0.20171127110929-2793f329227e
	github.com/lingochamp/xorm v0.6.4-0.20181203061557-28fcd64c4212
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/pkg/errors v0.9.1
)

require (
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-xorm/builder v0.3.4 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
)

replace github.com/go-xorm/xorm => github.com/lingochamp/xorm v0.6.4-0.20181203061557-28fcd64c4212
