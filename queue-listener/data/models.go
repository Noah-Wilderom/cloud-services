package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	Id         string     `json:"id"`
	Payload    JobPayload `json:"payload"`
	ReservedAt *time.Time `json:"reserved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type JobPayload struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
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
		var payloadData []byte // Temporary variable to hold payload data

		err := rows.Scan(
			&job.Id,
			&payloadData, // Scan the payload data into a []byte variable
			&job.ReservedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		// Unmarshal the payload data into the JobPayload field
		err = json.Unmarshal(payloadData, &job.Payload)
		if err != nil {
			log.Println("Error unmarshaling payload data", err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (jp *JobPayload) Scan(value interface{}) error {
	// Check if the value is nil
	if value == nil {
		*jp = JobPayload{} // Set the JobPayload to an empty value
		return nil
	}

	// Check if the value is of []uint8 type (common type for database blobs)
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid data type for JobPayload")
	}

	// Unmarshal the JSON-encoded bytes into the JobPayload struct
	err := json.Unmarshal(bytes, jp)
	if err != nil {
		return err
	}

	return nil
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
		var payloadData []byte // Temporary variable to hold payload data

		err := rows.Scan(
			&job.Id,
			&payloadData, // Scan the payload data into a []byte variable
			&job.ReservedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		// Parse the JSON-encoded payload data into the JobPayload struct
		var payload JobPayload
		err = json.Unmarshal(payloadData, &payload)
		if err != nil {
			log.Println("Error unmarshaling payload data", err)
			return nil, err
		}

		// Assign the parsed payload to the job
		job.Payload = payload

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

// Refresh refreshes the record from the database
func (u *Job) Refresh() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, payload, reserved_at, created_at, updated_at from jobs where id = $1`

	row := db.QueryRowContext(ctx, query, u.Id)

	var payloadData []byte // Temporary variable to hold payload data

	err := row.Scan(
		&u.Id,
		&payloadData, // Scan the payload data into a []byte variable
		&u.ReservedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Unmarshal the payload data into the JobPayload field
	err = json.Unmarshal(payloadData, &u.Payload)
	if err != nil {
		return err
	}

	return nil
}

func RefreshById(id string) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, payload, reserved_at, created_at, updated_at from jobs where id = $1`

	row := db.QueryRowContext(ctx, query, id)

	var j Job
	var payloadData []byte // Temporary variable to hold payload data

	err := row.Scan(
		&j.Id,
		&payloadData, // Scan the payload data into a []byte variable
		&j.ReservedAt,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal the payload data into the JobPayload field
	err = json.Unmarshal(payloadData, &j.Payload)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// Update updates one job in the database, using the information
// stored in the receiver u
func (j *Job) SetReserved() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update jobs set
		reserved_at = $1,
		where id = $2
	`

	_, err := db.ExecContext(ctx, stmt,
		time.Now(),
		j.Id,
	)

	if err != nil {
		return err
	}

	_ = j.Refresh()

	return nil
}

func SetReservedById(id string) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update jobs set
		reserved_at = $1
		where id = $2
	`

	_, err := db.ExecContext(ctx, stmt,
		time.Now(),
		id,
	)

	if err != nil {
		return &Job{}, err
	}

	return RefreshById(id)
}

// Delete deletes one job from the database, by Job.ID
func (j *Job) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from jobs where id = $1`

	_, err := db.ExecContext(ctx, stmt, j.Id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new job into the database, and returns the ID of the newly inserted row
func Insert(job Job) (*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID string
	stmt := `insert into jobs (payload, reserved_at, created_at, updated_at)
        values ($1, $2, $3, $4) returning id`

	// Marshal the JobPayload field into JSON before insertion
	payloadData, err := json.Marshal(job.Payload)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowContext(ctx, stmt,
		payloadData, // Insert the JSON-encoded payload data
		job.ReservedAt,
		job.CreatedAt,
		job.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return nil, err
	}

	job.Id = newID

	// No need to call Refresh() here, as the newly inserted data should already be in the struct

	return &job, nil
}
