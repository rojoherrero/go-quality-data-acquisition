package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var environment string
var pool *pgxpool.Pool
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "MANUFACTURING APP ", log.LstdFlags|log.Llongfile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var e error

	if pool, e = pgxpool.Connect(ctx, os.Getenv("POSTGRES_DSN")); e != nil {
		logger.Fatalln("Unable to connect to Postgras. Shutting down")
	}

}

func main() {
	fmt.Println("HELLO WORLD")
}
