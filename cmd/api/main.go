package main

import (
	"database/sql"
	"flag"
	"os"

	_ "github.com/lib/pq"
	"github.com/thomascastle/tarsk/internal/data"
	"github.com/thomascastle/tarsk/internal/messaging"
	"github.com/thomascastle/tarsk/internal/search"
	"github.com/thomascastle/tarsk/internal/structuredlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configuration struct {
	db struct {
		dsn string
	}
	env     string
	limiter struct {
		burst   int
		enabled bool
		rps     float64
	}
	port int
}

type application struct {
	config              configuration
	logger              *structuredlog.Logger
	messageBrokerage    *messaging.TaskMessageBrokerage
	repositories        data.Repositories
	search              data.Search
	taskIndexRepository data.TaskIndexRepository
}

func main() {
	var config configuration

	flag.StringVar(&config.db.dsn, "db-dsn", "", "Data Source Name")

	flag.StringVar(&config.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&config.limiter.burst, "limiter-burst", 4, "Maximum burst")
	flag.BoolVar(&config.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&config.limiter.rps, "limiter-rps", 2, "Maximum requests per second")

	flag.IntVar(&config.port, "port", 4000, "Port number the server is listening on")

	flag.Parse()

	logger := structuredlog.New(os.Stdout, structuredlog.LevelInfo)

	r_client, e := messaging.NewClient()
	if e != nil {
		logger.Fatal(e, nil)
	}

	logger.Info("messaging client created", nil)

	db, e := openDB(config)
	if e != nil {
		logger.Fatal(e, nil)
	}
	defer db.Close()

	logger.Info("database connection pool established", nil)

	s_client, e := search.NewClient()
	if e != nil {
		logger.Fatal(e, nil)
	}

	logger.Info("search client created", nil)

	db_GORM, e := openDB_GORM(config)
	if e != nil {
		logger.Fatal(e, nil)
	}

	app := &application{
		config:              config,
		logger:              logger,
		messageBrokerage:    messaging.NewTaskMessageBrokerage(r_client),
		repositories:        data.NewRepositories(db),
		search:              *data.NewSearch(s_client),
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

func openDB_GORM(config configuration) (*gorm.DB, error) {
	dsn := config.db.dsn
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if e != nil {
		return nil, e
	}

	return db, nil
}
