package db

import (
	"cligpt/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName     = "cligpt.db"
	folderName = ".cligpt"
)

func getDbPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(homedir, folderName, dbName)
}

func getDb() *sql.DB {
	var db *sql.DB
	filePath := getDbPath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Default().Println("Database file not found, creating one...")
		createDbFile()

		db, err = sql.Open("sqlite3", filePath)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		db, err = sql.Open("sqlite3", filePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db
}

func createDbFile() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	pathToFolder := filepath.Join(homedir, folderName)

	if _, err := os.Stat(pathToFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathToFolder, 0775)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Create(getDbPath())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	db, err := sql.Open("sqlite3", getDbPath())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Default().Println("Creating database table...")
	createQuery := "CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY, messages JSON, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"
	_, err = db.Exec(createQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database file created at: ", f.Name())
}

func GetLastTenSessions() []types.Session {
	db := getDb()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sessions ORDER BY updated_at DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	println(rows)

	var sessions []types.Session
	for rows.Next() {
		var id int
		var messages string
		var updated_at string

		err = rows.Scan(&id, &messages, &updated_at)
		if err != nil {
			log.Fatal(err)
		}
		
		var messagesArray []types.Message
		err = json.Unmarshal([]byte(messages), &messagesArray)
		if err != nil {
			log.Fatal(err)
		}

		sessions = append(sessions, types.Session{
			ID:       id,
			Messages: messagesArray,
		})

	}

	if len(sessions) == 0 {
		log.Fatal("No sessions found")
	}

	return sessions
}

func CreateSession(messages []types.Message) types.Session {
	db := getDb()
	defer db.Close()

	jsonMessages, err := json.Marshal(messages)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	session, err := db.Exec("INSERT INTO sessions (messages) VALUES (?)", jsonMessages)
	if err != nil {
		log.Fatal(err)
	}

	result, err := session.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return types.Session{
		ID:       int(result),
		Messages: []types.Message{},
	}
}

func UpdateSession(id int, messages []types.Message) {
	db := getDb()
	defer db.Close()

	jsonMessages, err := json.Marshal(messages)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	_, err = db.Exec("UPDATE sessions SET messages = ?, updated_at = ? WHERE id = ?", jsonMessages, time.Now().UTC(), id)
	if err != nil {
		log.Fatal(err)
	}
}
