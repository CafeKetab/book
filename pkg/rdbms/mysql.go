package rdbms

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type mysqlWrapper struct {
	*rdbms
}

func NewMysql(cfg *Config) (RDBMS, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("Error openning connection:\n%v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Error ping database:\n%v", err)
	}

	return &mysqlWrapper{&rdbms{db: db}}, nil
}

func (db *mysqlWrapper) Migrate(source string, direction MigrateDirection) error {
	var migrator func(m *migrate.Migrate) error

	switch direction {
	case MigrateUp:
		migrator = func(m *migrate.Migrate) error { return m.Up() }
	case MigrateDown:
		migrator = func(m *migrate.Migrate) error { return m.Down() }
	default:
		return errors.New("invalid migration direction has been given")
	}

	return db.migrate(source, migrator)
}

func (db *mysqlWrapper) migrate(source string, migrator func(*migrate.Migrate) error) error {
	instance, err := mysql.WithInstance(db.db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("Error creating migrate instance\n%v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(source, "mysql", instance)
	if err != nil {
		return fmt.Errorf("Error loading migration files\n%v", err)
	}

	if err := migrator(migration); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Error doing migrations\n%v", err)
	}

	return nil
}
