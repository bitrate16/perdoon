package tracks

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"perdoon/src/state"
	"perdoon/src/track"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	state *state.State
	db    *sql.DB
	lock  sync.Mutex

	// Reusaaable statement
	insertStatement *sql.Stmt
}

func NewSQLite(state *state.State) track.Track {
	return &SQLite{state: state}
}

func (s *SQLite) Open() error {
	dbDir := filepath.Dir(s.state.Config.Database.Path)

	// Ensure directory exists
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", s.state.Config.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	s.db = db

	// Create table if not exists
	_, err = s.db.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            session UUID,
            timestamp BIGINT,
            strategy TEXT,
            protocol TEXT,
            server_port INTEGER,
            client_ip TEXT,
            client_port INTEGER,
            chain_id UUID,
            event_type TEXT,
            payload_size BIGINT,
            payload_data BLOB
        )
    `, s.state.Config.Database.Table))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Cache prepared statement
	s.insertStatement, err = s.db.Prepare(fmt.Sprintf(`
        INSERT INTO %s (
            session,
			timestamp,
			strategy,
			protocol,
			server_port,
			client_ip,
			client_port,
			chain_id,
			event_type,
			payload_size,
			payload_data
        ) VALUES (
            ?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
        )
    `, s.state.Config.Database.Table))
	if err != nil {
		s.db.Close()
		return fmt.Errorf("failed to prepare statement: %s", err)
	}

	return nil
}

func (s *SQLite) Close() error {
	if s.db != nil {
		s.insertStatement.Close()
		return s.db.Close()
	}
	return nil
}

func (s *SQLite) Event(
	protocol string,
	server_port int,
	client_ip net.IP,
	client_port int,
	chain_id uuid.UUID,
	event_type string,
	payload_size int,
	payload_data []byte,
) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.db == nil {
		return fmt.Errorf("database closed")
	}

	timestamp := time.Now().UnixNano()

	var strategy string
	if protocol == "tcp" {
		strategy = s.state.Config.TCP.Response.Strategy
	} else {
		strategy = s.state.Config.UDP.Response.Strategy
	}

	_, err := s.insertStatement.Exec(
		s.state.Session,
		timestamp,
		strategy,
		protocol,
		server_port,
		client_ip.String(),
		client_port,
		chain_id,
		event_type,
		payload_size,
		payload_data,
	)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %s", err)
	}

	return nil
}
