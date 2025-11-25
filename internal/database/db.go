package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// Use environment variable or fallback to local path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Fallback for development
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		dbPath = filepath.Join(basepath, "..", "..", "movies.db")
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(0)

	// SQLite optimizations
	_, err = DB.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL; 
		PRAGMA cache_size = -64000;
		PRAGMA temp_store = memory;
		PRAGMA mmap_size = 268435456;
		PRAGMA busy_timeout = 5000;
	`)
	if err != nil {
		log.Printf("⚠️ Could not set SQLite optimizations: %v", err)
	}

	// INITIALIZE RESPONSE MANAGER
	InitResponseManager(DB)
	
	log.Println("✅ Database connected successfully")
	return nil
}