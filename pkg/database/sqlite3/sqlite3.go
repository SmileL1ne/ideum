package sqlite3

import "database/sql"

// OpenDB opens connection to the database using standard sql library
// with given Data Source Name (DSN)
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	// Enable foreign keys (they are disabled by default for backwards compatibility)
	query := "PRAGMA foreign_keys = ON;"
	_, err = db.Exec(query)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
