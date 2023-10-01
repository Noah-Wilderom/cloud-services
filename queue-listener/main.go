package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/queue-listener/data"
	"github.com/Noah-Wilderom/cloud-services/queue-listener/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

const (
	gRPCPort = 5003
)

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
	wg.Add(2) // We have two goroutines to wait for

	// Start your gRPC server in a goroutine
	go func() {
		defer wg.Done() // Decrement the WaitGroup when done
		app.gRPCListen()
	}()

	go func() {
		defer wg.Done() // Decrement the WaitGroup when done
		app.listenQueue()
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}

func (app *Config) listenQueue() {

	for {
		jobs, err := app.Models.Job.GetUnhandledJobs()
		if err != nil {
			log.Println(err)
		}

		if len(jobs) > 0 {
			for _, job := range jobs {
				_ = job.Refresh()
				if job.ReservedAt == nil {
					log.Printf("New job [%v] found sending to worker...", job.Id)
					err = app.SendToWorker(job)
					if err != nil {
						log.Println(err)
					}
				}
			}
		} else {
			log.Println("No new jobs found...")
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func (app *Config) SendToWorker(job *data.Job) error {
	conn, err := grpc.Dial(
		"queue-worker:5002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := queue.NewQueueWorkerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println(string(job.Payload.Data))

	jobPayload := &queue.JobPayload{
		Service: job.Payload.Service,
		Data:    job.Payload.Data,
	}

	jobResponse, err := c.HandleJob(ctx, &queue.JobRequest{
		Job: &queue.Job{
			Id:         job.Id,
			Payload:    jobPayload,
			ReservedAt: &timestamppb.Timestamp{},
			UpdatedAt:  timestamppb.New(job.UpdatedAt),
			CreatedAt:  timestamppb.New(job.CreatedAt),
		},
	})
	if err != nil {
		return err
	}

	if jobResponse.Error {
		return errors.New(jobResponse.ErrorPayload)
	}

	return nil
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
