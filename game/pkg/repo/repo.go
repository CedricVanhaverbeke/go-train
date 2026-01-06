package repo

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"overlay/pkg/gpx"

	_ "modernc.org/sqlite"
)

type GPXRepo struct {
	db *sql.DB
}

type GPXRecord struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Data      string    `json:"data"` // GPX data stored as JSON string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGPXRepo(dbPath string) (*GPXRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &GPXRepo{db: db}

	if err := repo.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return repo, nil
}

func (r *GPXRepo) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS gpx_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		data TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	_, err := r.db.Exec(query)
	return err
}

func (r *GPXRepo) Create(name string, gpxData gpx.Gpx) (*GPXRecord, error) {
	dataXML, err := xml.Marshal(gpxData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GPX data: %w", err)
	}

	now := time.Now()
	query := `
	INSERT INTO gpx_files (name, data, created_at, updated_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, name, string(dataXML), now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert GPX record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return &GPXRecord{
		ID:        id,
		Name:      name,
		Data:      string(dataXML),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *GPXRepo) Get(id int64) (*GPXRecord, error) {
	query := `
	SELECT id, name, data, created_at, updated_at
	FROM gpx_files
	WHERE id = ?
	`

	record := &GPXRecord{}
	err := r.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.Name,
		&record.Data,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GPX record with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get GPX record: %w", err)
	}

	return record, nil
}

func (r *GPXRepo) GetGPX(id int64) (*gpx.Gpx, error) {
	record, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	var gpxData gpx.Gpx
	err = json.Unmarshal([]byte(record.Data), &gpxData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GPX data: %w", err)
	}

	return &gpxData, nil
}

func (r *GPXRepo) GetAll() ([]*GPXRecord, error) {
	query := `
	SELECT id, name, data, created_at, updated_at
	FROM gpx_files
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query GPX records: %w", err)
	}
	defer rows.Close()

	var records []*GPXRecord
	for rows.Next() {
		record := &GPXRecord{}
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Data,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan GPX record: %w", err)
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating GPX records: %w", err)
	}

	return records, nil
}

func (r *GPXRepo) Update(id int64, name string, gpxData gpx.Gpx) (*GPXRecord, error) {
	dataJSON, err := json.Marshal(gpxData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GPX data: %w", err)
	}

	now := time.Now()
	query := `
	UPDATE gpx_files
	SET name = ?, data = ?, updated_at = ?
	WHERE id = ?
	`

	_, err = r.db.Exec(query, name, string(dataJSON), now, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update GPX record: %w", err)
	}

	return r.Get(id)
}

func (r *GPXRepo) Delete(id int64) error {
	query := `DELETE FROM gpx_files WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete GPX record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("GPX record with ID %d not found", id)
	}

	return nil
}

func (r *GPXRepo) Close() error {
	return r.db.Close()
}
