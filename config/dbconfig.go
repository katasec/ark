package config

type DbConfig struct {
	DriverName     string // DriverName is database driver database/sql.Open uses
	DataSourceName string // DataSourceName is the data source name database/sql.Open uses
}
