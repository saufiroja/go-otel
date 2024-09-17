package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/saufiroja/go-otel/auth-service/config"
	"github.com/saufiroja/go-otel/auth-service/pkg/logging"
	"sync"
)

type PostgresManager interface {
	Connection() *sql.DB
	StartTransaction() (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
	CloseConnection() error
}

type Postgres struct {
	db *sql.DB
}

var (
	postgresInstance *Postgres
	once             sync.Once
)

func NewPostgres(conf *config.AppConfig, logger logging.Logger) PostgresManager {
	once.Do(func() {
		user := conf.Postgres.User
		password := conf.Postgres.Pass
		dbHost := conf.Postgres.Host
		dbPort := conf.Postgres.Port
		dbName := conf.Postgres.Name
		dbSslMode := conf.Postgres.SSL

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, user, password, dbName, dbSslMode)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			logger.LogPanic(fmt.Sprintf("Error opening databases: %v", err))
		}

		if err := db.Ping(); err != nil {
			logger.LogPanic(fmt.Sprintf("Error connecting to databases: %v", err))
		}

		logger.LogInfo("Database connected")
		postgresInstance = &Postgres{db: db}
	})

	return postgresInstance
}

func (p *Postgres) Connection() *sql.DB {
	return p.db
}

func (p *Postgres) StartTransaction() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *Postgres) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (p *Postgres) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func (p *Postgres) CloseConnection() error {
	return p.db.Close()
}
