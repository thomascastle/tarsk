package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("server started at :4000")
	if e := server.ListenAndServe(); e != nil {
		log.Fatal("failed to serve at :4000")
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
