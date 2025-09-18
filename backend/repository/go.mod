module github.com/kenkonno/gantt-chart-proto/backend/repository

go 1.23.0

require (
	github.com/jackc/pgx/v5 v5.7.5
	github.com/kenkonno/gantt-chart-proto/backend/models v0.0.1
	github.com/samber/lo v1.37.0
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.24.2
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/kenkonno/gantt-chart-proto/backend/models v0.0.1 => ../models
