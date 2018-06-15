package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func reportError(err error, w http.ResponseWriter, r *http.Request, userString string, code int) {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	fmt.Printf("Error [%d] %s: %s (%s)\n", code, r.URL.String(), userString, errString)
	http.Error(w, userString, code)
}

func main() {
	postgresAddr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		os.Getenv("PG_HOST"), os.Getenv("PG_DBNAME"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_SSL_MODE"))
	database, err := sql.Open("postgres", postgresAddr)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()

	isProduction := flag.Bool("production", false, "whether to run in production mode")
	flag.Parse()

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	gzFileServer := gziphandler.GzipHandler(http.FileServer(http.Dir("./static")))
	http.Handle("/", gzFileServer)

	if *isProduction {
		log.Fatal(http.ListenAndServe(":80", redirectToHTTPS(http.DefaultServeMux)))
	} else {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}
}

func redirectToHTTPS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHttps := r.Header.Get("X-Forwarded-Proto") == "https"
		if !isHttps && r.URL.Path != "/healthcheck" {
			// Adapted from https://gist.github.com/d-schmidt/587ceec34ce1334a5e60
			target := "https://" + os.Getenv("HOSTNAME") + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
