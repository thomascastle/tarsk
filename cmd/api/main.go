package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/thomascastle/tarsk/internal/data"
)

type configuration struct {
	db struct {
		dsn string
	}
	port int
}

type application struct {
	config configuration
	models data.Models
}

func main() {
	var config configuration

	flag.StringVar(&config.db.dsn, "db-dsn", "", "Data Source Name")
	flag.IntVar(&config.port, "port", 4000, "Port number the server is listening on")

	flag.Parse()

	db, e := openDB(config)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()

	log.Println("database connection pool established")

	app := &application{
		config: config,
		models: data.NewModels(db),
	}

	e = app.serve()
	if e != nil {
		log.Fatal(e)
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
