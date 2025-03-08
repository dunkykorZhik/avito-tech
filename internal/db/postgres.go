package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/dunkykorZhik/avito-tech/config"
)

type postgresDB struct {
	Db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (Database, error) {
	p := &postgresDB{}
	var err error

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.DbName,
		cfg.Db.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Db.MaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err

	}

	p.Db = db

	return p, err

}

func (p *postgresDB) GetDb() *sql.DB {
	if p.Db == nil {
		panic("ahahahaha again")
	}
	return p.Db
}
