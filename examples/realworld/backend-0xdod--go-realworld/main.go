package main

import (
	"fmt"
	"log"
	"os"

	"github.com/0xdod/go-realworld/postgres"
	"github.com/0xdod/go-realworld/server"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	port  string
	dbURI string
}

func main() {
	cfg := envConfig()

	fmt.Printf("%s\n", cfg.dbURI)
	db, err := postgres.Open(cfg.dbURI)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Run(cfg.port))
}

func envConfig() config {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		panic("PORT not provided")
	}

	dbURI, ok := os.LookupEnv("POSTGRESQL_URL")

	if !ok {
		panic("POSTGRESQL_URL not provided")
	}

	return config{port, dbURI}
}
