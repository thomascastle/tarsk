package main

import (
	"database/sql"
	"flag"
	"os"

	_ "github.com/lib/pq"
	"github.com/thomascastle/tarsk/internal/data"
	"github.com/thomascastle/tarsk/internal/structuredlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configuration struct {
	db struct {
		dsn string
	}
	env  string
	port int
}

type application struct {
	config              configuration
	logger              *structuredlog.Logger
	repositories        data.Repositories
	taskIndexRepository data.TaskIndexRepository
}

func main() {
	var config configuration

	flag.StringVar(&config.db.dsn, "db-dsn", "", "Data Source Name")

	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&config.port, "port", 4000, "Port number the server is listening on")

	flag.Parse()

	logger := structuredlog.New(os.Stdout, structuredlog.LevelInfo)

	db, e := openDB(config)
	if e != nil {
		logger.Fatal(e, nil)
	}
	defer db.Close()

	logger.Info("database connection pool established", nil)

	db_GORM, e := openDB_GORM()
	if e != nil {
		logger.Fatal(e, nil)
	}

	app := &application{
		config:              config,
		logger:              logger,
		repositories:        data.NewRepositories(db),
		taskIndexRepository: data.NewTaskIndexRepository(db_GORM),
	}

	e = app.serve()
	if e != nil {
		logger.Fatal(e, nil)
	}
}

func openDB(config configuration) (*sql.DB, error) {
	db, e := sql.Open("postgres", config.db.dsn)
	if e != nil {
		return nil, e
	}

	e = db.Ping()
	if e != nil {
		return nil, e
	}

	return db, nil
}

func openDB_GORM() (*gorm.DB, error) {
	dsn := "postgres://thomas:securesecret@localhost/tarsk?sslmode=disable"
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if e != nil {
		return nil, e
	}

	return db, nil
}
