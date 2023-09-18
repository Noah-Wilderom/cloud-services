package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Job: Job{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	Job Job
}

// Job is the structure which holds one job from the database.
type Job struct {
	ID         string     `json:"id"`
	Payload    string     `json:"payload"`
	ReservedAt *time.Time `json:"reserved_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// GetAll returns a slice of all jobs, sorted by last name
func (u *Job) GetAll() ([]*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, payload, reserved_at, created_at, updated_at
	from jobs order by created_at desc`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Job

	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID,
			&job.Payload,
			&job.ReservedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (u *Job) GetUnhandledJobs() ([]*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, payload, reserved_at, created_at, updated_at
	from jobs where reserved_at is null order by created_at desc`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*Job

	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.ID,
			&job.Payload,
			&job.ReservedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

// Refresh refreshes the record from the database
func (u *Job) Refresh() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, payload, reserved_at, created_at, updated_at from jobs where id = $1`

	row := db.QueryRowContext(ctx, query, u.ID)

	err := row.Scan(
		u.ID,
		u.Payload,
		u.ReservedAt,
		u.CreatedAt,
		u.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// Update updates one job in the database, using the information
// stored in the receiver u
func (j *Job) SetReserved() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update jobs set
		reserved_at= $1,
		where id = $2
	`

	_, err := db.ExecContext(ctx, stmt,
		time.Now(),
		j.ID,
	)

	if err != nil {
		return err
	}

	_ = j.Refresh()

	return nil
}

// Delete deletes one job from the database, by Job.ID
func (j *Job) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from jobs where id = $1`

	_, err := db.ExecContext(ctx, stmt, j.ID)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new job into the database, and returns the ID of the newly inserted row
func (j *Job) Insert(job Job) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID string
	stmt := `insert into jobs (payload, reserved_at, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := db.QueryRowContext(ctx, stmt,
		job.Payload,
		time.Now(),
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return &job, err
	}

	job.ID = newID

	_ = job.Refresh()

	return &job, nil
}
