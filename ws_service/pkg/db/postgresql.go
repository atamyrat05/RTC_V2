package package_psql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Options struct {
	Host          string
	Port          string
	Database      string
	Username      string
	Password      string
	PgPoolMaxConn int
}

func NewClient(ctx context.Context, options Options) (Client, error) {
	log.Println("new client options")
	log.Println("🔔 Host: ", options.Host)
	log.Println("🔔 Port: ", options.Port)
	log.Println("🔔 Database: ", options.Database)
	log.Println("🔔 Username: ", options.Username)
	log.Println("🔔 Password: ", "*******")
	log.Println("🔔 PgPoolMaxConn: ", options.PgPoolMaxConn)

	connPool, err := pgxpool.NewWithConfig(ctx, getConfig(options))
	if err != nil {
		return nil, errors.New("🚫 Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		return nil, errors.New("🚫 Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Conn().Ping(ctx)
	if err != nil {
		return nil, errors.New("🚫 Could not ping database")
	}

	log.Println("✅ postgresql connected success")

	// *sql.DB oluşturmak için pgx stdlib paketini kullanıyoruz
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
	)
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	stdDB := stdlib.OpenDB(*config)

	if err := stdDB.PingContext(ctx); err != nil {
		return nil, errors.New("🚫 Could not ping stdDB")
	}

	return &clientImpl{
		pool:  connPool,
		stdDB: stdDB,
	}, nil
}

func getConfig(options Options) *pgxpool.Config {
	databaseURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
	)

	log.Println("🔔 database url: ", databaseURL)

	dbConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Println("🚫 Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = int32(options.PgPoolMaxConn)
	dbConfig.MinConns = int32(0)
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = time.Minute * 30
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5

	return dbConfig
}
