package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// Entry represents a paste entry
type Entry struct {
	ID           int64
	URL          string
	Content      string
	EditCode     string
	ModifyCode   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ViewCount    int
	Deleted      bool
	CreatedByIP  string
	LastEditedIP string
}

// User represents a registered user
type User struct {
	ID        int64
	Username  string
	Password  string // hashed
	CreatedAt time.Time
	IsAdmin   bool
}

// NewDB creates a new database connection and initializes tables
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return database, nil
}

// initTables creates the necessary database tables
func (db *DB) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		edit_code TEXT NOT NULL,
		modify_code TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		view_count INTEGER DEFAULT 0,
		deleted BOOLEAN DEFAULT 0,
		created_by_ip TEXT,
		last_edited_ip TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_entries_url ON entries(url);
	CREATE INDEX IF NOT EXISTS idx_entries_deleted ON entries(deleted);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`

	_, err := db.Exec(schema)
	return err
}

// CreateEntry inserts a new entry into the database
func (db *DB) CreateEntry(entry *Entry) error {
	query := `
		INSERT INTO entries (url, content, edit_code, modify_code, created_by_ip)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, entry.URL, entry.Content, entry.EditCode, entry.ModifyCode, entry.CreatedByIP)
	if err != nil {
		return err
	}

	entry.ID, err = result.LastInsertId()
	return err
}

// GetEntryByURL retrieves an entry by its URL
func (db *DB) GetEntryByURL(url string) (*Entry, error) {
	entry := &Entry{}
	query := `
		SELECT id, url, content, edit_code, modify_code, created_at, updated_at, 
		       view_count, deleted, created_by_ip, last_edited_ip
		FROM entries
		WHERE url = ? AND deleted = 0
	`
	err := db.QueryRow(query, url).Scan(
		&entry.ID, &entry.URL, &entry.Content, &entry.EditCode, &entry.ModifyCode,
		&entry.CreatedAt, &entry.UpdatedAt, &entry.ViewCount, &entry.Deleted,
		&entry.CreatedByIP, &entry.LastEditedIP,
	)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// UpdateEntry updates an existing entry
func (db *DB) UpdateEntry(entry *Entry) error {
	query := `
		UPDATE entries
		SET content = ?, edit_code = ?, modify_code = ?, updated_at = CURRENT_TIMESTAMP, 
		    last_edited_ip = ?, url = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, entry.Content, entry.EditCode, entry.ModifyCode, 
		entry.LastEditedIP, entry.URL, entry.ID)
	return err
}

// DeleteEntry marks an entry as deleted
func (db *DB) DeleteEntry(id int64) error {
	query := `UPDATE entries SET deleted = 1 WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// HardDeleteEntry permanently deletes an entry (admin only)
func (db *DB) HardDeleteEntry(id int64) error {
	query := `DELETE FROM entries WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// IncrementViewCount increases the view count for an entry
func (db *DB) IncrementViewCount(id int64) error {
	query := `UPDATE entries SET view_count = view_count + 1 WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// URLExists checks if a URL is already taken
func (db *DB) URLExists(url string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM entries WHERE url = ? AND deleted = 0`
	err := db.QueryRow(query, url).Scan(&count)
	return count > 0, err
}

// CreateUser creates a new user
func (db *DB) CreateUser(username, hashedPassword string) error {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := db.Exec(query, username, hashedPassword)
	return err
}

// GetUserByUsername retrieves a user by username
func (db *DB) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, created_at, is_admin FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllEntries retrieves all entries (for admin)
func (db *DB) GetAllEntries(limit, offset int) ([]*Entry, error) {
	query := `
		SELECT id, url, content, edit_code, modify_code, created_at, updated_at, 
		       view_count, deleted, created_by_ip, last_edited_ip
		FROM entries
		WHERE deleted = 0
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		entry := &Entry{}
		err := rows.Scan(
			&entry.ID, &entry.URL, &entry.Content, &entry.EditCode, &entry.ModifyCode,
			&entry.CreatedAt, &entry.UpdatedAt, &entry.ViewCount, &entry.Deleted,
			&entry.CreatedByIP, &entry.LastEditedIP,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// SearchEntries searches entries by URL or content
func (db *DB) SearchEntries(query string, limit int) ([]*Entry, error) {
	sqlQuery := `
		SELECT id, url, content, edit_code, modify_code, created_at, updated_at, 
		       view_count, deleted, created_by_ip, last_edited_ip
		FROM entries
		WHERE deleted = 0 AND (url LIKE ? OR content LIKE ?)
		ORDER BY created_at DESC
		LIMIT ?
	`
	searchPattern := "%" + query + "%"
	rows, err := db.Query(sqlQuery, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*Entry
	for rows.Next() {
		entry := &Entry{}
		err := rows.Scan(
			&entry.ID, &entry.URL, &entry.Content, &entry.EditCode, &entry.ModifyCode,
			&entry.CreatedAt, &entry.UpdatedAt, &entry.ViewCount, &entry.Deleted,
			&entry.CreatedByIP, &entry.LastEditedIP,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetStatistics returns general statistics
func (db *DB) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalEntries, totalViews, totalUsers int
	
	err := db.QueryRow(`SELECT COUNT(*) FROM entries WHERE deleted = 0`).Scan(&totalEntries)
	if err != nil {
		return nil, err
	}
	
	err = db.QueryRow(`SELECT COALESCE(SUM(view_count), 0) FROM entries WHERE deleted = 0`).Scan(&totalViews)
	if err != nil {
		return nil, err
	}
	
	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	if err != nil {
		return nil, err
	}

	stats["total_entries"] = totalEntries
	stats["total_views"] = totalViews
	stats["total_users"] = totalUsers

	return stats, nil
}
