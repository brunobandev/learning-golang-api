package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"vue-api/internal/driver"
)

// config is the type for all application configuration.
type config struct {
	port int
}

// application is the type for all data we want to share with the
// various parts of our application. We will share this information in most
// cases by using this type as the receiver for functions
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	db       *driver.DB
}

// main is the main entry point for our application
func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dns := "host=localhost port=5433 user=postgres password=password dbname=vueapi sslmode=disable timezone=UTC connect_timeout=5"
	db, err := driver.ConnectPostgres(dns)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       db,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

// serve starts the web server
func (app *application) serve() error {
	app.infoLog.Println("API Listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
