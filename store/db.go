// Package store maneja la persistencia de perfiles analizados.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"distroanalyzer/profile"
)

// Store representa un almacén persistente de perfiles.
type Store interface {
	Save(ctx context.Context, p *profile.Profile) error
	GetByUsername(ctx context.Context, username string) (*profile.Profile, error)
	List(ctx context.Context, limit, offset int) ([]*profile.Profile, error)
	Delete(ctx context.Context, username string) error
}

// SQLiteStore implementa Store usando SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore crea un store basado en SQLite.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Crear tablas si no existen
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return store, nil
}

// migrate crea las tablas necesarias.
func (s *SQLiteStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS profiles (
		username TEXT PRIMARY KEY,
		source TEXT NOT NULL,
		raw_data TEXT,
		signals TEXT NOT NULL,
		result TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_created_at ON profiles(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_source ON profiles(source);
		`

		_, err := s.db.Exec(query)
		return err
}

// Save guarda o actualiza un perfil.
func (s *SQLiteStore) Save(ctx context.Context, p *profile.Profile) error {
	// Serializar structs complejos a JSON
	rawDataJSON, err := json.Marshal(p.RawData)
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %w", err)
	}

	signalsJSON, err := json.Marshal(p.Signals)
		if err != nil {
			return fmt.Errorf("failed to marshal signals: %w", err)
		}

		resultJSON, err := json.Marshal(p.Result)
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		now := time.Now()

		query := `
		INSERT INTO profiles (username, source, raw_data, signals, result, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
		source = excluded.source,
		raw_data = excluded.raw_data,
		signals = excluded.signals,
		result = excluded.result,
		updated_at = excluded.updated_at
		`

		_, err = s.db.ExecContext(ctx, query,
					  p.Username,
					  p.Source,
					  string(rawDataJSON),
					  string(signalsJSON),
					  string(resultJSON),
					  p.CreatedAt,
					  now,
		)

		return err
}

// GetByUsername obtiene un perfil por username.
func (s *SQLiteStore) GetByUsername(ctx context.Context, username string) (*profile.Profile, error) {
	query := `
	SELECT username, source, raw_data, signals, result, created_at
	FROM profiles
	WHERE username = ?
	`

	var p profile.Profile
	var rawDataJSON, signalsJSON, resultJSON string

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&p.Username,
		&p.Source,
		&rawDataJSON,
		&signalsJSON,
		&resultJSON,
		&p.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Deserializar JSON
	if err := json.Unmarshal([]byte(rawDataJSON), &p.RawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw data: %w", err)
	}

	if err := json.Unmarshal([]byte(signalsJSON), &p.Signals); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signals: %w", err)
	}

	if err := json.Unmarshal([]byte(resultJSON), &p.Result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &p, nil
}

// List obtiene perfiles paginados.
func (s *SQLiteStore) List(ctx context.Context, limit, offset int) ([]*profile.Profile, error) {
	query := `
	SELECT username, source, raw_data, signals, result, created_at
	FROM profiles
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*profile.Profile

	for rows.Next() {
		var p profile.Profile
		var rawDataJSON, signalsJSON, resultJSON string

		err := rows.Scan(
			&p.Username,
		   &p.Source,
		   &rawDataJSON,
		   &signalsJSON,
		   &resultJSON,
		   &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Deserializar JSON
		if err := json.Unmarshal([]byte(rawDataJSON), &p.RawData); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(signalsJSON), &p.Signals); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(resultJSON), &p.Result); err != nil {
			return nil, err
		}

		profiles = append(profiles, &p)
	}

	return profiles, rows.Err()
}

// Delete elimina un perfil.
func (s *SQLiteStore) Delete(ctx context.Context, username string) error {
	query := `DELETE FROM profiles WHERE username = ?`
	_, err := s.db.ExecContext(ctx, query, username)
	return err
}

// Close cierra la conexión a la base de datos.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
