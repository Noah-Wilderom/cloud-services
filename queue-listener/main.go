package main

import (
	"database/sql"
	"github.com/Noah-Wilderom/cloud-services/queue-listener/data"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var (
	counts int64
)

func main() {

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	var wg sync.WaitGroup

	// Increment the WaitGroup to indicate a goroutine is starting
	wg.Add(1)

	// Start your gRPC server in a goroutine
	go func() {
		defer wg.Done() // Decrement the WaitGroup when done
		app.listenQueue()
	}()

	// Wait for the gRPC server to start
	wg.Wait()

}

func (app *Config) listenQueue() {

	for {
		jobs, err := app.Models.Job.GetUnhandledJobs()
		if err != nil {
			log.Println(err)
		}

		if len(jobs) < 1 {
			log.Println("No new jobs found...")
			return
		}

		for _, job := range jobs {
			_ = job.Refresh()
			if job.ReservedAt == nil {
				log.Printf("New job [%v] found sending to worker...", job.ID)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func openDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDBConnection(dsn)
		if err != nil {
			log.Println()
			log.Println("Postgres not yet ready...", err)
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)

			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
