package mariadb

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"github.com/voxtmault/mentoring/library-project/pkg/storage"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func validateConfig(cfg *config.MariaDBConfig) error {
	if cfg.DBUser == "" {
		return eris.New("db username is empty")
	}
	if cfg.DBPassword == "" {
		return eris.New("db password is empty")
	}
	if cfg.DBHost == "" || cfg.DBPort == "" {
		return eris.New("invalid db address and or port")
	}
	if cfg.DBName == "" {
		return eris.New("invalid db name")
	}

	return nil
}

func Init(cfg *config.MariaDBConfig) (*sql.DB, error) {
	slog.Debug("init mariadb connection")

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	dsn := mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		AllowNativePasswords: cfg.AllowNativePasswords,
		Addr:                 cfg.DBHost + ":" + cfg.DBPort,
		DBName:               cfg.DBName,
		TLSConfig:            cfg.TLSConfig,
		MultiStatements:      cfg.MultiStatements,
		Net:                  "tcp",
		Params: map[string]string{
			"charset":   "utf8",
			"parseTime": "true",
		},
	}

	con, err := sql.Open(storage.MariaDBDriver, dsn.FormatDSN())
	if err != nil {
		return nil, eris.Wrap(err, "failed to open connection")
	}

	con.SetMaxOpenConns(int(cfg.MaxOpenConns))
	con.SetMaxIdleConns(int(cfg.MaxIdleConns))
	con.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifetime))

	if err = con.Ping(); err != nil {
		return nil, eris.Wrap(err, "Error verifying database connection")
	}

	slog.Info("mariadb connection established", "host", cfg.DBHost)
	return con, nil
}

func InitORM(con *sql.DB) (*gorm.DB, error) {
	if con == nil {
		return nil, eris.New("mariadb connection is nil")
	}

	ormCon, err := gorm.Open(
		gormMysql.New(
			gormMysql.Config{
				Conn: con,
			},
		), &gorm.Config{
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to open gorm connection")
	}

	return ormCon, nil
}

func Close(con *sql.DB) error {
	if err := con.Close(); err != nil {
		return eris.Wrap(err, "failed to close connection")
	}

	return nil
}
